/*
* 标签处理
*
* 标签以 hash 的格式存储在 redis 里，
* 并且以在 mysql 中的主键作为存在 redis 里的 field 名，
* 以次来映射 redis 里标签和 mysql 里的标签的关系
*
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
)

// TagForm 定义提交的tag表单
type TagForm struct {
	Color    string `form:"color" json:"color" binding:"required"`
	TagTitle string `form:"tag_title" json:"tag_title" binding:"required"`
}

// GetAllTagHandler 获取全部标签
func GetAllTagHandler(c *gin.Context) {
	// 从 redis 里获取全部标签信息
	tags, err := models.RedisClient.HVals("tags").Result()

	// 如果从 redis 里读取失败或者不存在就从数据库里读取
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Get all tags from redis failed")

		// 从数据库获取全部分类名
		tags, err := models.GetAllTag()
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    http.StatusText(http.StatusInternalServerError),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"message":    "Query tags failed",
				"statusCode": http.StatusInternalServerError,
			}).Info("Get all tags failed")

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"tags":       tags,
		})

		// 同步到 redis 里
		for _, tag := range *tags {
			b, err := models.RedisClient.HSet("tags", string(tag.ID), tag).Result()
			if !b || err != nil {
				log.WithFields(log.Fields{
					"errorMsg": err,
					"tag":      tag,
				}).Info("Sync to to redis failed")
			}
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"tags":       tags,
	})

	log.WithFields(log.Fields{
		"tags":       tags,
		"statusCode": http.StatusOK,
	}).Info("Get all tags success")
}

// AddTagHandler 增加标签
func AddTagHandler(c *gin.Context) {
	var tagVals TagForm

	// 管理员鉴权
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 检测提交的表单
	err = c.ShouldBindWith(&tagVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Miss color or tag_title",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Add tag failed")

		return
	}

	// 存储进数据库
	lastID, err := models.AddTag(tagVals.Color, tagVals.TagTitle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store tag to database failed",
			"statusCode": http.StatusBadRequest,
		}).Info("Add tag failed")

		return
	}

	c.JSON(200, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Add tag success",
	})

	log.WithFields(log.Fields{
		"tag":        tagVals,
		"statusCode": http.StatusOK,
	}).Info("Add tag to mysql success")

	// 同步更新到 redis 里
	b, err := models.RedisClient.HSet("tags", string(lastID), models.Tag{ID: uint(lastID), Color: tagVals.Color, TagTitle: tagVals.TagTitle}).Result()
	if !b || err != nil {
		log.WithFields(log.Fields{
			"id":  lastID,
			"tag": tagVals,
		}).Info("Store tag to redis failed")
	}

	log.WithFields(log.Fields{
		"id":       lastID,
		"tagColor": tagVals.Color,
		"tagTitle": tagVals.TagTitle,
	}).Info("Sync tag to redis success")
}

// 就检查是否重复提交已存在标签
func checkRepeatTag(key string, field string) (bool, error) {
	// 首先开始从 redis 里查询
	tags, err := models.RedisClient.HVals(key).Result()
	if err != nil {
		// 从数据库里查询
		tags, err := models.GetAllTag()
		if err != nil {
			return false, err
		}

		for _, tag := range *tags {
			if tag.TagTitle == field {
				return true, nil
			}
		}

		return false, nil
	}

	for _, tag := range tags {
		log.WithFields(log.Fields{
			"tag": tag,
		}).Info("tagxxxxxxx")
		// if tag.TagTitle == field {
		// 	return true, nil
		// }
	}

	return false, nil
}
