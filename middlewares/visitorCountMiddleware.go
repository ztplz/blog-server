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
		var isAddSuccess bool
		var visitorCount uint
		var isGetAllCount bool
		var isAddRedisSuccess bool

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

		// 尝试五次插入数据库，知道成功为止
		for i := 0; i < 5; i++ {
			err := models.CountVistor()
			if err == nil {
				isAddSuccess = true
				log.Info("Count visitor success")

				break
			}
		}

		if !isAddSuccess {
			log.Info("Count visitor failed")
		}

		// 数据库更新访问人数成功，则同步到 redis, 用 for 循环保证更高的成功几率
		if isAddSuccess {
			for i := 0; i < 5; i++ {
				count, err := models.GetAllVisitorCount()
				if err == nil {
					isGetAllCount = true
					visitorCount = count
				}
			}

			if !isGetAllCount {
				return
			}

			// 把访问统计人数同步到 redis 里
			for i := 0; i < 5; i++ {
				err := models.RedisClient.Set("visitor_count", string(visitorCount), 0).Err()
				if err == nil {
					isAddRedisSuccess = true

					break
				}
			}

			// 如果同步到 redis 失败，打印日志
			if !isAddRedisSuccess {
				log.Info("Sync visitor count to redis success")
			}
		}
	}
}

// ClearIPPool 清空 IP 池
func ClearIPPool() {
	IPPool = IPPool[:0]
	log.Info("clear success")
}
