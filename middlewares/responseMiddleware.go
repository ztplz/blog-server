/*
* 生成响应 模块
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

// package middlewares

// import (
// 	"github.com/gin-gonic/gin"
// 	log "github.com/sirupsen/logrus"
// )

// type LogContent struct {
// 	message interface{}
// }

// ResponseMiddleware 返回写入中间件
// func ResponseMiddleware(c *gin.Context, requestType string, statusCode int, message string, logInfo string, logContent interface{}, data interface{}) {
// 	switch requestType {
// 	case "GET":
// 		// 设置为 json 格式
// 		c.JSON(statusCode, gin.H{
// 			"statusCode": statusCode,
// 			"message":    message,
// 			"data":       data,
// 		})

		// 禁止请求继续执行
// 		c.AbortWithStatus(statusCode)

// 		// 打印相关信息
// 		log.WithFields(log.Fields{
// 			"message":    logContent,
// 			"statusCode": statusCode,
// 		}).Info(logInfo)
// 	}
	
// 	case "POST":
// 		// 设置为 json 格式
// 		c.JSON(statusCode, gin.H{
// 			"statusCode": statusCode,
// 			"message":    message,
// 		})

// 		// 禁止请求继续执行
// 		c.AbortWithStatus(statusCode)

// 		// 打印相关信息
// 		log.WithFields(log.Fields{
// 			"message":    logContent,
// 			"statusCode": statusCode,
// 		}).Info(logInfo)
	
	

// }
