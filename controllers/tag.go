package controllers

import (
	"errors"
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
	tags, err := models.GetAllTag()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "query failed",
		})
		c.AbortWithError(401, errors.New("query failed"))

		return
	}

	c.JSON(200, gin.H{
		"tags": tags,
	})
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
	err = models.AddTag(tagVals.Color, tagVals.TagTitle)
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
	}).Info("Add tag success")
}
