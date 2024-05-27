package svc

import (
	"demo-gateway/common/stores/xgorm"
	"demo-gateway/rpc/study/internal/config"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DshDb  *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	dshDb := xgorm.NewGorm(xgorm.GormConfig{
		Dsn:             c.DshMysql.Dsn,
		Tracing:         c.DshMysql.Tracing,
		MaxIdle:         c.DshMysql.MaxIdle,
		MaxOpen:         c.DshMysql.MaxOpen,
		ConnMaxIdleTime: c.DshMysql.ConnMaxIdleTime,
	}, c.Log)
	return &ServiceContext{
		Config: c,
		DshDb:  dshDb,
	}
}
