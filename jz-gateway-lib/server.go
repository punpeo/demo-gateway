package gateway

import (
	"context"
	"demo-gateway/jz-gateway-lib/internal"
	"fmt"
	"github.com/punpeo/punpeo-lib/rest/result"
	"github.com/punpeo/punpeo-lib/rest/xerr"
	"github.com/punpeo/punpeo-lib/utils/jzcrypto"
	"net/http"
	"strings"

	"github.com/fullstorydev/grpcurl"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc/codes"
)

type (
	// Server is a gateway server.
	Server struct {
		*rest.Server
		upstreams []Upstream
		//jz-gateway 添加 http.Request 拦截的异常捕获
		processHeader func(http.Header, *http.Request) []string
		dialer        func(conf zrpc.RpcClientConf) zrpc.Client
		/*new*/
		Config *GatewayConf
		plugin *PluginManager
	}

	// Option defines the method to customize Server.
	Option func(svr *Server)
)

// MustNewServer creates a new gateway server.
func MustNewServer(c *GatewayConf, opts ...Option) *Server {
	svr := &Server{
		upstreams: c.Upstreams,
		Server:    rest.MustNewServer(c.RestConf),
		/*new*/
		Config: c,
		plugin: NewPluginManager(),
	}
	for _, opt := range opts {
		opt(svr)
	}

	return svr
}

// Register 注册插件
/*new*/
func (s *Server) Register(p Plugin) {
	s.plugin.Register(p)
}

// Start starts the gateway server.
func (s *Server) Start() {
	/*new*/
	LoadRouteMap(s.Config)
	logx.Must(s.build())
	s.Server.Start()
}

// Stop stops the gateway server.
func (s *Server) Stop() {
	s.Server.Stop()
}

func (s *Server) build() error {
	if err := s.ensureUpstreamNames(); err != nil {
		return err
	}

	return mr.MapReduceVoid(func(source chan<- Upstream) {
		for _, up := range s.upstreams {
			source <- up
		}
	}, func(up Upstream, writer mr.Writer[rest.Route], cancel func(error)) {
		var cli zrpc.Client
		if s.dialer != nil {
			cli = s.dialer(up.Grpc)
		} else {
			cli = zrpc.MustNewClient(up.Grpc)
		}

		source, err := s.createDescriptorSource(cli, up)
		if err != nil {
			cancel(fmt.Errorf("%s: %w", up.Name, err))
			return
		}

		methods, err := internal.GetMethods(source)
		if err != nil {
			cancel(fmt.Errorf("%s: %w", up.Name, err))
			return
		}

		resolver := grpcurl.AnyResolverFromDescriptorSource(source)
		for _, m := range methods {
			if len(m.HttpMethod) > 0 && len(m.HttpPath) > 0 {
				route := rest.Route{
					Method:  m.HttpMethod,
					Path:    m.HttpPath,
					Handler: s.buildHandler(source, resolver, cli, m.RpcPath, false), /*new*/
				}

				// 设置中间件
				route = s.plugin.WrapMiddleware(&route)
				writer.Write(route)
			}
		}

		methodSet := make(map[string]struct{})
		for _, m := range methods {
			methodSet[m.RpcPath] = struct{}{}
		}

		for _, m := range up.Mappings {
			// 加载路由对应插件
			s.plugin.LoadRouteMapping(&up, &m) /*new*/

			if _, ok := methodSet[m.RpcPath]; !ok {
				cancel(fmt.Errorf("%s: rpc method %s not found", up.Name, m.RpcPath))
				return
			}

			// OrigName配置
			origName := up.OrigName
			if m.OrigName != nil {
				origName = *m.OrigName
			}

			route := rest.Route{
				Method:  strings.ToUpper(m.Method),
				Path:    m.Path,
				Handler: s.buildHandler(source, resolver, cli, m.RpcPath, origName), /*new*/
			}

			// 设置中间件
			route = s.plugin.WrapMiddleware(&route) /*new*/
			writer.Write(route)                     /*new*/
		}
	}, func(pipe <-chan rest.Route, cancel func(error)) {
		for route := range pipe {
			s.Server.AddRoute(route)
		}
	})
}

func (s *Server) buildHandler(source grpcurl.DescriptorSource, resolver jsonpb.AnyResolver,
	cli zrpc.Client, rpcPath string, origName bool) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		parser, err := internal.NewRequestParser(r, resolver)
		if err != nil {
			//jz-gateway 调整返回值
			httpx.WriteJson(w, http.StatusOK, &result.ResponseSuccessBean{Code: xerr.SERVER_COMMON_ERROR, Msg: err.Error(), Data: "请求参数解析错误"})
			return
		}
		w.Header().Set(httpx.ContentType, httpx.JsonContentType)

		// 设置RPC事件处理器
		//handler := internal.NewEventHandler(w, resolver) //返回格式未做处理
		handler := s.plugin.GetRpcHandler(w, r, resolver, origName) //处理返回格式

		if err := grpcurl.InvokeRPC(r.Context(), source, cli.Conn(), rpcPath, s.prepareMetadata(r.Header, r),
			handler, parser.Next); err != nil {
			//jz-gateway 调整返回值
			logx.Errorf("rpc调用失败,%+v", err.Error())
			httpx.WriteJson(w, http.StatusOK, &result.ResponseSuccessBean{Code: xerr.SERVER_COMMON_ERROR, Msg: err.Error()})
			return
		}

		st := handler.Status
		if st.Code() != codes.OK {
			//jz-gateway 调整返回值
			// if handler.XStatusCode != 0 { //自定义code
			// 	httpRespMsg := st.Message()
			// 	if handler.XErrorMessage != "" {
			// 		httpRespMsg = handler.XErrorMessage
			// 	}
			// 	httpx.WriteJson(w, http.StatusOK, &result.ResponseSuccessBean{Code: handler.XStatusCode, Msg: httpRespMsg})
			// 	return
			// }
			logx.Errorf("rpc响应失败,%+v", st.Err())
			httpx.WriteJson(w, http.StatusOK, &result.ResponseSuccessBean{Code: xerr.SERVER_COMMON_ERROR, Msg: st.Message()})
			return
		}
	}
}

func (s *Server) createDescriptorSource(cli zrpc.Client, up Upstream) (grpcurl.DescriptorSource, error) {
	var source grpcurl.DescriptorSource
	var err error

	if len(up.ProtoSets) > 0 {
		source, err = grpcurl.DescriptorSourceFromProtoSets(up.ProtoSets...)
		if err != nil {
			return nil, err
		}
	} else {
		client := grpcreflect.NewClientAuto(context.Background(), cli.Conn())
		source = grpcurl.DescriptorSourceFromServer(context.Background(), client)
	}

	return source, nil
}

func (s *Server) ensureUpstreamNames() error {
	for i := 0; i < len(s.upstreams); i++ {
		target, err := s.upstreams[i].Grpc.BuildTarget()
		if err != nil {
			return err
		}

		s.upstreams[i].Name = target
	}

	return nil
}

func (s *Server) prepareMetadata(header http.Header, req *http.Request) []string {
	vals := internal.ProcessHeaders(header)
	if s.processHeader != nil {
		vals = append(vals, s.processHeader(header, req)...)
	}

	return vals
}

// WithHeaderProcessor sets a processor to process request headers.
// The returned headers are used as metadata to invoke the RPC.
func WithHeaderProcessor(processHeader func(http.Header, *http.Request) []string) func(*Server) {
	return func(s *Server) {
		s.processHeader = processHeader
	}
}

// withDialer sets a dialer to create a gRPC client.
func withDialer(dialer func(conf zrpc.RpcClientConf) zrpc.Client) func(*Server) {
	return func(s *Server) {
		s.dialer = dialer
	}
}

// LoadRouteMap 加载配置转map
func LoadRouteMap(c *GatewayConf) {
	AuthCheckMapping := make(map[string]map[string]bool)
	VerifyFuncControlMapping := make(map[string]map[string]bool)
	UpstreamsRouteMap := make(map[string]map[string]RouteMapping)
	for _, upstream := range c.Upstreams {
		for _, mapping := range upstream.Mappings {
			if _, ok := AuthCheckMapping[strings.ToLower(mapping.Method)]; !ok {
				AuthCheckMapping[strings.ToLower(mapping.Method)] = map[string]bool{strings.ToLower(mapping.Path): mapping.AuthCheck}
			}
			AuthCheckMapping[strings.ToLower(mapping.Method)][strings.ToLower(mapping.Path)] = mapping.AuthCheck
			if _, ok := VerifyFuncControlMapping[strings.ToLower(mapping.Method)]; !ok {
				VerifyFuncControlMapping[strings.ToLower(mapping.Method)] = map[string]bool{strings.ToLower(mapping.Path): mapping.AuthCheck}
			}
			VerifyFuncControlMapping[strings.ToLower(mapping.Method)][strings.ToLower(mapping.Path)] = mapping.VerifyFuncControl
			if _, ok := UpstreamsRouteMap[strings.ToLower(mapping.Method)]; !ok {
				UpstreamsRouteMap[strings.ToLower(mapping.Method)] = make(map[string]RouteMapping)
			}
			UpstreamsRouteMap[strings.ToLower(mapping.Method)][strings.ToLower(mapping.Path)] = mapping

		}
	}
	c.AuthCheckMapping = AuthCheckMapping
	c.VerifyFuncControlMapping = VerifyFuncControlMapping
	c.UpstreamsRouteMap = UpstreamsRouteMap
}

type
// GatewayConf is the configuration for gateway.
GatewayConf struct {
	rest.RestConf
	Upstreams []Upstream
	//管理后台相关权限rpc服务
	AccessControlRpc zrpc.RpcClientConf
	//是否校验强制登录map
	AuthCheckMapping map[string]map[string]bool `json:",optional"`
	//是否校验功能权限map
	VerifyFuncControlMapping map[string]map[string]bool `json:",optional"`
	//路由配置Map
	UpstreamsRouteMap map[string]map[string]RouteMapping
	//app、h5的security_key签名
	Safe Safe
	//php内部调用sign
	SafeList    []jzcrypto.TripleDesConf
	SignKey     string
	SignKeyList []string
}

// RouteMapping is a mapping between a gateway route and an upstream rpc method.
type RouteMapping struct {
	// Method is the HTTP method, like GET, POST, PUT, DELETE.
	Method string
	// Path is the HTTP path.
	Path string
	// RpcPath is the gRPC rpc method, with format of package.service/method
	RpcPath string
	// AuthCheck token检查，默认为检查
	AuthCheck bool `json:",optional,default=true"`
	// Plugins 单一路由的插件，将完全覆盖全局插件
	Plugins []string `json:",optional"`
	// OrigName 单一路由控制 是否启用OriginName  未配置则使用 Upstream.OrigName
	OrigName *bool `json:",optional"`
	// VerifyFuncControl 功能权限检查，默认为不检查
	VerifyFuncControl bool        `json:",optional,default=false"`
	UriDispatch       UriDispatch `json:",optional"`
}

// Upstream is the configuration for an upstream.
type Upstream struct {
	// Name is the name of the upstream.
	Name string `json:",optional"`
	// Grpc is the target of the upstream.
	Grpc zrpc.RpcClientConf
	// ProtoSets is the file list of proto set, like [hello.pb].
	// if your proto file import another proto file, you need to write multi-file slice,
	// like [hello.pb, common.pb].
	ProtoSets []string `json:",optional"`
	// Mappings is the mapping between gateway routes and Upstream rpc methods.
	// Keep it blank if annotations are added in rpc methods.
	Mappings []RouteMapping `json:",optional"`
	// Plugins 全局插件
	Plugins []string `json:",optional"`
	// OrigName  是否启用OriginName 默认不开启
	OrigName bool `json:",optional,default=false"`
}

type Safe struct {
	Key string
	Iv  string
}

type UriDispatch struct {
	//调度方案 0-直连go服务  1-用户灰度方案(新旧接口，不同服务) 2-兜底双请求校验 3-直连php服务 4-内部版本灰度(同服务接口，不同版本) 5-灰度＋兜底
	DispatchRule int8 `json:",optional,default=0"`
	//兜底优先 [缺省]0-php  2-go
	Priority int8 `json:",optional,default=0"`
	//灰度方案 1-用户取模 userId % GrayDivisor 2-配置模式 (指定用户ID) 3-配置模式 (指定ip)
	GrayScheme int8 `json:",optional"`
	//灰度比例 GrayScheme == 1 用户模: [0,1,2,3]
	GrayRate []int8 `json:",optional"`
	//灰度取模 除数
	GrayDivisor int16 `json:",optional,default=100"`
	//直连php host
	DirectHost string `json:",optional"`
	//直连php地址
	DirectPath string `json:",optional"`
	//配置模式 文件路径
	GrayConfigPath string `json:",optional"`
	//内部版本灰度 [缺省]0-php  1-go, 增加灰度头 canary:1
	DispatchServer int8 `json:",optional"`
}
