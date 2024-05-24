package main

import (
	"flag"
	"fmt"
	"google.golang.org/grpc/reflection"

	"demo-gateway/rpc/study/internal/config"
	studyServer "demo-gateway/rpc/study/internal/server/study"
	"demo-gateway/rpc/study/internal/svc"
	"demo-gateway/rpc/study/study"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

var configFile = flag.String("f", "etc/study.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		study.RegisterStudyServer(grpcServer, studyServer.NewStudyServer(ctx))

		//if c.Mode == service.DevMode || c.Mode == service.TestMode {
		//	reflection.Register(grpcServer)
		//}
		reflection.Register(grpcServer)
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
