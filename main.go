/*
* blog 的 main 程序，程序从这里开始运行
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/models"
	"github.com/ztplz/blog-server/router"
)

func main() {
	// 日志设置
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// 初始化数据库连接
	models.InitialDB()

	// 初始化管理员账户
	models.InitialAdmin()

	// 初始化路由
	router.InitialRouter()

	defer models.DB.Close()
}
