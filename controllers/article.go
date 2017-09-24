package controllers

import (
	"net/http"
	"strconv"
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

// GetAllArticlesHandler 获取文章列表
func GetAllArticlesHandler(c *gin.Context) {
	// 获取查询参数
	limitString := c.DefaultQuery("limit", "10")
	pageString := c.DefaultQuery("page", "1")

	// 转换 limitString, page 为 int 类型
	limit, lerr := strconv.ParseInt(limitString, 10, 32)
	page, perr := strconv.ParseInt(pageString, 10, 32)
	if lerr != nil || perr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Parameter not incorrect",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsgLimit": lerr,
			"errorMsgPage":  perr,
			"limit":         limitString,
			"pageString":    pageString,
			"statusCode":    http.StatusBadRequest,
		}).Info("Get articles failed")

		return
	}

	// 判断是否大于 0
	if limit <= 0 || page <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Parameter must above 0",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Parameter must above 0",
			"limit":      limit,
			"pageString": page,
			"statusCode": http.StatusBadRequest,
		}).Info("Get articles failed")

		return
	}

	// 查询 博文
	articles, err := models.GetArticleByPage(limit, page)
	for _, article := range *articles {
		*&article.TagList = strings.Split(*&article.TagList, "_")
	}
	articles.TagList = strings.Split(articles.TagList, "_")

	c.JSON(http.StatusBadRequest, gin.H{
		"statusCode": http.StatusOK,
		"data":       articles,
	})
}

// AddArticleHandler 增加文章
func AddArticleHandler(c *gin.Context) {
	var articleVals ArticleForm

	// 验证是否有权限增加文章
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 从表单中提取文章标题
	err = c.ShouldBindWith(&articleVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Article form miss some field",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Article form incorrect")

		return
	}

	// 判断标题是否规定长度
	if len(articleVals.ArticleTitle) > models.ArticleTitleLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Article title too long",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Article title too long",
			"statusCode": http.StatusBadRequest,
		}).Info("Article form incorrect")

		return
	}

	// 预览内容是否规定长度
	if len(articleVals.ArticlePreviewText) > models.ArticlePreviewTextLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Article preview text too long",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Article preview text too long",
			"statusCode": http.StatusBadRequest,
		}).Info("Article form incorrect")

		return
	}

	// 判断标签是否规定数量
	if len(articleVals.Tags) > 3 {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Too many tags",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Too many tags",
			"statusCode": http.StatusBadRequest,
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

	// 存储博文数据进数据库
	err = models.AddArticle(article)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
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

	// 打印相关日志
	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
		"title":      articleVals.ArticleTitle,
		"create_at":  time.Now().Format("2006-01-02 15:04:05"),
	}).Info("Add article success")
}
