package controllers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
)

// ArticleForm 文章表单
type ArticleForm struct {
	ArticleTitle       string   `form:"article_title" json:"article_title" binding:"required"`
	ArticlePreviewText string   `form:"article_previewtext" json:"article_previewtext" binding:"required"`
	ArticleContent     string   `form:"article_content" json:"article_content" binding:"required"`
	Category           string   `form:"category" json:"category" binding:"required"`
	Tags               []string `form:"tags" json:"tags" binding:"required"`
}

// AddArticleHandler 增加文章
func AddArticleHandler(c *gin.Context) {
	articleVals := &ArticleForm{}

	// 验证是否有权限增加文章
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 从表单中提取文章标题
	if c.ShouldBindWith(articleVals, binding.JSON) != nil {
		c.JSON(400, gin.H{
			"message":    err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err.Error(),
			"statusCode": http.StatusBadRequest,
		}).Info("Article form incorrect")

		return
	}

	// 判断标签是否规定数量
	if len(articleVals.Tags) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    "Too many tags",
			"statusCode": http.StatusBadRequest,
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Too many tags",
			"statusCode": http.StatusInternalServerError,
		}).Info("Add article failed")

		return
	}

	// 把传到后台的标签数组转化成字符串
	tags := strings.Join(articleVals.Tags, "_")

	// 实例化一个 Article 结构体
	article := &models.Article{
		CreateAt:           time.Now().Format("2006-01-02 15:04:05"),
		UpdateAt:           time.Now().Format("2006-01-02 15:04:05"),
		VisitCount:         0,
		ReplyCount:         0,
		ArticleTitle:       articleVals.ArticleTitle,
		ArticlePreviewText: articleVals.ArticlePreviewText,
		ArticleContent:     articleVals.ArticleContent,
		Top:                false,
		Category:           articleVals.Category,
		TagList:            tags,
	}

	err = models.AddArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":    http.StatusText(http.StatusInternalServerError),
			"statusCode": http.StatusInternalServerError,
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store article to database failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Add article failed")

		return
	}

	c.JSON(200, gin.H{
		"success": "true",
		"message": "Add article success",
	})

}
