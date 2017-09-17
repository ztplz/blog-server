package main

import (
	"github.com/ztplz/blog-server/router"
	"github.com/ztplz/blog-server/models"
)

func main() {
	// 初始化数据库连接
	models.InitialDB()

	// 初始化管理员账户
	models.InitialAdmin()

	// 初始化路由
	router.InitialRouter()
		
	defer models.DB.Close()
}