package studylogic

import (
	"context"

	"demo-gateway/rpc/study/internal/svc"
	"demo-gateway/rpc/study/study"

	"github.com/zeromicro/go-zero/core/logx"
)

type HelloLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHelloLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HelloLogic {
	return &HelloLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *HelloLogic) Hello(in *study.HelloRequest) (*study.HelloResponse, error) {
	// todo: add your logic here and delete this line

	return &study.HelloResponse{
		Message: "hello,欢迎使用go-zero网关",
	}, nil
}
