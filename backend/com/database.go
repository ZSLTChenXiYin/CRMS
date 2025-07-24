package com

import (
	"log"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDatabase() error {
	// 启动数据库
	db, err := gorm.Open(mysql.Open(Config.DSN), &gorm.Config{})
	if err != nil {
		return err
	}
	Database = db

	log.Println("[INFO] MySQL连接成功")

	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := Redis.Ping(Context).Result()
	if err != nil {
		return err
	}

	log.Println("[INFO] Redis连接成功:", pong)

	return nil
}
