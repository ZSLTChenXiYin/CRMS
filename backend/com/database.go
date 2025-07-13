package com

import (
	"log"

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

	log.Println("[INFO] 数据库连接成功")

	return nil
}
