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
	"github.com/ztplz/blog-server/middlewares"
	"gopkg.in/robfig/cron.v2"
)

func main() {
	// 定时器设置
	c := cron.New()
	c.AddFunc("0 0 0 * * *", middlewares.AddTodayVisitorCount)
	c.Start()
	// 日志设置
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	// 初始化连接redis
	models.InitialRedis()
	defer models.RedisClient.Close()

	// 初始化数据库连接
	models.InitialDB()
	defer models.DB.Close()

	// 初始化管理员账户
	models.InitialAdmin()

	// 初始化路由
	router.InitialRouter()

	// 定时任务
	// select{}
	
}
