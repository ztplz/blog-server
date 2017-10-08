/*
* 访客统计中间件，00:00 - 24:00内，计算pv
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package middlewares

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/models"
)

// IPPool ip池，00：00清空一次
var IPPool []string

// CountVisitorMiddleware  统计访问人数中间件
func CountVisitorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var isRepeat bool
		// var isAddSuccess bool
		// var visitorCount uint
		// var isGetAllCount bool
		// var isAddRedisSuccess bool

		ip := c.ClientIP()
		for _, value := range IPPool {
			if ip == value {
				isRepeat = true

				break
			}
		}

		// 如果ip跟每天ip池里的不重复
		if !isRepeat {
			IPPool = append(IPPool, ip)
		}
    
		// 今日 IP 池数量
		// ipCount := len(IPPool)

		// 每日访问人数同步到redis里
		// err := models.RedisClient.Set("today_visitor_count", ipCount, 0).Err()
		// if err != nil {
		// 	log.WithFields(log.Fields{
		// 		"errorMsg": err,
		// 	}).Info("Store today visitor count to redis failed")
		// }
	}
}

//AddTodayVisitorCount 添加昨日访问人数到数据库
func AddYesterdayVisitorCount() {
	var isAddToDatabaseSuccess bool

	// 清楚 redis 里历史访问人数
	err := models.RedisClient.Del("all_visitor_count").Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Delete all visitor count from redis failed")
	}

	count := uint(len(IPPool))
	// 历史访问人数比较重要，最多尝试五次确保成功
	for i := 0; i < 5; i++ {
		err := models.CountVistor(count)
		if err == nil {
			isAddToDatabaseSuccess = true

			break
		}
	}

	// 数据库存储失败
	if !isAddToDatabaseSuccess {
		log.Info("Insert yestday visitor count to database failed")
	}

	// 清空 IP 池
	clearIPPool()
}

//  清空 IP 池
func clearIPPool() {
	IPPool = IPPool[:0]
	log.Info("clear success")
}
