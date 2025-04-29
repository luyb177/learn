package dao

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"learn/check_status/config"
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

	// 增加重试机制
	var err error
	for i := 0; i < 3; i++ {
		// 测试连接
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err = rdb.Ping(ctx).Result()
		if err == nil {
			cancel()        // 取消之前的上下文
			return rdb, nil // 连接成功，返回客户端
		}
		log.Printf("Redis连接失败，正在重试... (%d/3)\n", i+1)
		time.Sleep(5 * time.Second)
		cancel() // 取消之前的上下文
	}
	return nil, fmt.Errorf("redis连接失败")
}

func MySqlDB(cfg *config.MysqlConfig) (*gorm.DB, error) {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // 输出到标准输出
		logger.Config{
			SlowThreshold: time.Second, // 慢查询阈值，超过此时间的查询将被记录
			LogLevel:      logger.Info, // 记录信息级别（Info、Warn、Error）
			Colorful:      true,        // 输出带颜色
		},
	)
	var db *gorm.DB
	var sqlDB *sql.DB
	var err error
	for i := 0; i < 5; i++ {
		db, err = gorm.Open(mysql.Open(cfg.Addr), &gorm.Config{
			Logger: newLogger,
		})

		if err == nil {
			// 获取数据库连接实例并设置连接池参数
			sqlDB, err = db.DB()
			if err != nil {
				return nil, fmt.Errorf("failed to get generic database object: %v", err)
			}

			// 设置连接池参数
			sqlDB.SetMaxOpenConns(100) // 设置最大打开连接数
			sqlDB.SetMaxIdleConns(10)  // 设置最大空闲连接数

			return db, nil
		}
		log.Printf("MySQL连接失败，正在重试... (%d/5)\n", i+1)
		time.Sleep(5 * time.Second)
	}
	if err != nil {
		panic("MySQL连接失败")
	}
	return nil, fmt.Errorf("MySQL连接失败")
}
