package plugins

import (
	gateway "demo-gateway/jz-gateway-lib"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	"google.golang.org/grpc/metadata"
	"net/http"
	"strconv"
)

// PluginEmpty 空插件，为屏蔽某一路由的插件使用
type PluginEmpty struct {
	gateway.BasicRpcHandler
	gw *gateway.Server
}

func NewPluginEmpty() *PluginEmpty {
	return &PluginEmpty{}
}

func (p *PluginEmpty) Name() string {
	return "empty"
}

func (p *PluginEmpty) Middleware() rest.Middleware {
	hdl := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			if r.Method == http.MethodGet {
				logx.WithContext(ctx).Info("empty")
			}
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
	return rest.ToMiddleware(hdl)
	//return nil
}
func (p *PluginEmpty) OnReceiveResponse(respJson string, md metadata.MD, _ http.ResponseWriter) string {
	respCode, respMsg := 1000, "成功"

	xStatusCodeArr := md.Get("X-Status-Code")
	if len(xStatusCodeArr) > 0 {
		codeInt64, _ := strconv.ParseInt(xStatusCodeArr[0], 10, 64)
		respCode = int(codeInt64)
	}

	xErrorMessage := md.Get("X-Error-Message")
	if len(xErrorMessage) > 0 {
		respMsg = xErrorMessage[0]
	}

	xData := md.Get("X-Data")
	if len(xData) > 0 {
		respJson = xData[0]
	}

	return fmt.Sprintf("{\"code\":%d, \"msg\": \"%s\", \"data\":%s}", respCode, respMsg, respJson)
}
