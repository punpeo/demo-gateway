package auth

import (
	"context"
	"encoding/json"
	"fmt"
	xerr2 "github.com/punpeo/punpeo-lib/rest/xerr"
	"github.com/spf13/cast"
	"time"

	commonTypes "demo-gateway/common/types"
)

// 获取账号数据权限信息
func GetAccountDataAuth(ctx context.Context) *commonTypes.AccountDataAuth {
	data := ctx.Value("account_data_auth")
	dataAuth := &commonTypes.AccountDataAuth{
		Type: 0,
	}
	if data == nil {
		return dataAuth
	}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return dataAuth
	}
	json.Unmarshal(jsonStr, dataAuth)

	return dataAuth
}

func GetCreators(ctx context.Context) (string, error) {
	var creators string
	auth := GetAccountDataAuth(ctx)
	if auth.Type == 0 {
		return "", xerr2.NewErrCode(xerr2.MISSED_DATA_PERMISSIONS_ERROR)
	}
	if auth.Type == 2 {
		if len(auth.List) > 0 {
			for v, k := range auth.List {
				if v == 0 {
					creators = fmt.Sprintf("%v", k.Id)
				} else {
					creators = fmt.Sprintf("%v,%v", creators, k.Id)
				}
			}
		} else {
			creators = "0"
		}
	}
	return creators, nil
}

// 获取登录user_id
func GetLoginUserId(ctx context.Context) int64 {
	var userId int64
	var expireTime int64
	expireTime = cast.ToInt64(ctx.Value("expire_time"))
	if expireTime > time.Now().Unix() {
		userId = cast.ToInt64(ctx.Value("user_id"))
	}
	return userId
}

// 是否登录
func IsLogin(ctx context.Context) error {
	userId := GetLoginUserId(ctx)
	if userId <= 0 {
		return xerr2.NewErrCode(xerr2.LOGIN_EXPIRE_ERROR)
	}
	return nil
}
