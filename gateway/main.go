package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/gateway"
)

var configFile = flag.String("f", "etc/study-rpc-gateway.yaml", "the config file")

func main() {
	flag.Parse()
	var c gateway.GatewayConf
	logx.DisableStat()
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c)
	//loadPlugins(gw)
	defer gw.Stop()
	gw.Start()

}
