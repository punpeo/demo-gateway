package studylogic

import (
	"context"
	model "demo-gateway/common/model/dushuhui"
	"errors"

	"demo-gateway/rpc/study/internal/svc"
	"demo-gateway/rpc/study/study"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserInfoLogic {
	return &UserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserInfoLogic) UserInfo(in *study.UserInfoRequest) (*study.UserInfoResponse, error) {
	var resp *study.UserInfoResponse
	resp = &study.UserInfoResponse{}
	if in.Id == 0 {
		return nil, errors.New("id不能为空")
	}
	DshUserWeixinModel := model.DshUserWeixin{}
	userInfo, err := DshUserWeixinModel.GetInfo(l.svcCtx.DshDb, l.ctx, in.Id)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if userInfo != nil {
		resp = &study.UserInfoResponse{
			Id: userInfo.Id,
		}
	}
	return resp, nil
}
