package main

import (
	gateway "demo-gateway/jz-gateway-lib"
	"demo-gateway/jz-gateway-lib/plugins"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/study-rpc-gateway.yaml", "the config file")

func main() {
	flag.Parse()
	var c gateway.GatewayConf
	logx.DisableStat()
	conf.MustLoad(*configFile, &c) //加载配置文件, 并解析到c中
	gw := gateway.MustNewServer(&c)
	loadPlugins(gw)
	defer gw.Stop()
	gw.Start()

}

// loadPlugins 加载插件
func loadPlugins(gw *gateway.Server) {
	gw.Register(plugins.NewPluginJzAuth(gw.Config))
	gw.Register(plugins.NewPluginEmpty())
	gw.Register(plugins.NewPluginHls())
	gw.Register(plugins.NewPluginUriDispatch(gw.Config))
	gw.Register(plugins.NewPluginCustom())
}
