package config

import (
	"demo-gateway/common/stores/xgorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DshMysql xgorm.GormConfig
}
