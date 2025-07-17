# 原生部署 #

## 环境要求 ##
* [Go 1.23+](https://go.dev/dl/)
* [MySQL 8.0+](https://www.mysql.com/downloads/)
* [Nginx 1.24+](https://nginx.org/en/download.html)

## 步骤 ##
1. 创建MySQL用户和数据库

2. 修改API服务配置文件（backend/config.json），保证MySQL配置正确

3. 运行API服务
```bash
cd backend && go run main.go
```

4. 通过Nginx挂载frontend目录搭建Web服务
