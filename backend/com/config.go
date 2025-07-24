package com

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type config struct {
	ReleaseMode   bool   `json:"release_mode"`
	AutoMigrate   bool   `json:"auto_migrate"`
	DSN           string `json:"dsn"`
	RedisAddress  string `json:"redis_address"`
	RedisPassword string `json:"redis_password"`
	JwtSecret     []byte `json:"jwt_secret"`
	Port          int    `json:"port"`
}

func InitConfig() error {
	// 读取配置文件
	conf_file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer conf_file.Close()

	conf_data, err := io.ReadAll(conf_file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(conf_data, Config)
	if err != nil {
		return err
	}

	log.Println("[INFO] 配置数据加载成功")

	return nil
}
