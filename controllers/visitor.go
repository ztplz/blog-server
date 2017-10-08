package controllers

import (
	"net/http"

	"github.com/ztplz/blog-server/middlewares"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/models"
)

// GetVisitCount 获取访问人数信息
func GetVisitCount(c *gin.Context) {
	count, err := models.RedisClient.Get("all_visitor_count").Result()
	// 从 redis 里读取失败
	if err != nil {
		log.WithFields(log.Fields{
			"message": err,
		}).Info("Get visitor count failed")

		// 从数据库读取
		count, err := models.GetAllVisitorCount()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    http.StatusText(http.StatusInternalServerError),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"message":    "Query all visitor count from database failed",
				"statusCode": http.StatusInternalServerError,
			}).Info("Get all visitor count failed")

			return
		}

		// 如果从数据库读取成功就同步到 redis 里
		err = models.RedisClient.Set("all_visitor_count", string(count), 0).Err()
		if err != nil {
			log.WithFields(log.Fields{
				"message": err,
			}).Info("Sync all visitor count to redis failed")
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode":          http.StatusOK,
		"message":             "获取访客数据成功",
		"all_visitor_count":   count,
		"today_visitor_count": len(middlewares.IPPool),
	})
}
