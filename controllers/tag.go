package controllers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
)

// 获取全部标签
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

// 增加标签
func AddTagHandler(c *gin.Context) {
	_, err := middlewares.AdminAuthMiddleware(c)

	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		c.AbortWithError(401, errors.New("auth failed"))

		return
	}

	// 要修改的分类名
	color := c.Query("color")
	log.Println(color)
	// 替换原来的分类名
	title := c.Query("title")
	log.Println(title)

	err = models.AddTag(color, title)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "failed",
		})
		c.AbortWithError(401, errors.New("Add failed"))

		return
	}

	c.JSON(200, gin.H{
		"message": "Add success",
	})
}
