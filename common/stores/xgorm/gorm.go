package xgorm

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

type GormConfig struct {
	Dsn             string
	Tracing         bool
	MaxIdle         int
	MaxOpen         int
	ConnMaxIdleTime int
}

func NewGorm(gormConfig GormConfig, config logx.LogConf) *gorm.DB {
	var (
		filePath  string
		newLogger logger.Interface
	)
	if config.Path == "" {
		newLogger = logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	} else {
		filePath = fmt.Sprintf("%v/sql-%v.log", config.Path, time.Now().Format("2006-01-02"))
		file, _ := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
		newLogger = logger.New(
			log.New(file, "\r\n", log.LstdFlags), // 文件
			//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: false,       // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
	}
	db, err := gorm.Open(mysql.Open(gormConfig.Dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(gormConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(gormConfig.MaxOpen)
	sqlDB.SetConnMaxIdleTime(time.Duration(int64(gormConfig.ConnMaxIdleTime)) * time.Second)

	return db
}
