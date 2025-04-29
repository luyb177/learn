package dao

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"learn/Grab_seat/config"
	"log"
	"os"
	"time"
)

func RedisDB(cfg *config.RedisConfig) (*redis.Client, error) {
	// 创建Redis客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,     // redis地址
		Password: cfg.Password, // Redis认证密码(可选)
		DB:       cfg.DB,       // 选择的数据库
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect redis: %v", err)
	}

	return rdb, nil
}

func MySQLDB(cfg *config.MysqlConfig) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值，超过此时间的查询将被记录
			LogLevel:      logger.Info, // 记录信息级别（Info、Warn、Error）
			Colorful:      true,        // 输出带颜色
		},
	)
	db, err := gorm.Open(mysql.Open(cfg.Addr), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic("failed to connect database")
	}
	// 获取数据库连接实例并设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get generic database object: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(10)  // 设置最大空闲连接数

	return db, nil
}
