package controllers

import (
	"encoding/json"
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
	Category           uint     `form:"category" json:"category" binding:"required"`
	Tags               []string `form:"tags" json:"tags" binding:"required"`
}

// ArticleRes 定义博文响应的结构体
type ArticleRes struct {
	ID                 uint            `json:"id"`
	CreateAt           string          `json:"creat_at"`
	UpdateAt           string          `json:"update_at"`
	VisitCount         uint            `json:"visit_count"`
	ReplyCount         uint            `json:"reply_count"`
	ArticleTitle       string          `json:"article_title"`
	ArticlePreviewText string          `json:"article_previewtext"`
	ArticleContent     string          `json:"article_content"`
	Top                bool            `json:"top"`
	Category           models.Category `json:"category"`
	TagList            []models.Tag    `json:"tag_list"`
}

// GetAllArticlesHandler 获取文章列表
func GetAllArticlesHandler(c *gin.Context) {
	var articlesRes []ArticleRes
	// 获取查询参数
	limitString := c.DefaultQuery("limit", "10")
	pageString := c.DefaultQuery("page", "1")

	// 转换 limitString, page 为 int64 类型
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

	// 首先从 redis 里查询博文
	articles, err := models.RedisClient.LRange("articles", (page-1)*limit, page*limit).Result()
	marticles := new([]models.Article)
	marticle := new(models.Article)
	for _, article := range articles {
		err := json.Unmarshal([]byte(article), &marticle)
		if err != nil {
			log.WithFields(log.Fields{
				"errorMsg":   err,
				"limit":      limit,
				"pageString": page,
			}).Info("Get articles from redis failed")

			break
		}

		*marticles = append(*marticles, *marticle)
	}

	// 用长度判断反序列化过程中有没有错误，是否该从mysql拉取数据
	if len(*marticles) == int(limit) {
		for _, value := range *marticles {

		}
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusOk,
			"articles":   *marticles,
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

	// 首先从 redis 里查询所有的分类 ID\

	keys, err := models.RedisClient.HKeys("category").Result()
	ukeys := new([]uint)
	for _, value := range keys {
		uvalue := value
		ukeys = append(ukeys)
	}

	categories, err := models.GetAllCategory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Query category failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Get articles failed")

		return
	}

	// 查询所有标签
	tags, err := models.GetAllTag()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Query tag failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Get articles failed")

		return
	}

	// 查询博文
	articles, err := models.GetArticleByPage(limit, page)
	for _, article := range *articles {
		log.Info("categories", categories)
		category := categoryForRes(article.Category, &categories)
		tagList, err := tagForRes(article.TagList, tags)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    http.StatusText(http.StatusInternalServerError),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"errorMsg":   err,
				"statusCode": http.StatusInternalServerError,
			}).Info("Get articles failed")
		}
		articlesRes = append(articlesRes, ArticleRes{
			ID:                 article.ID,
			CreateAt:           article.CreateAt,
			UpdateAt:           article.UpdateAt,
			VisitCount:         article.VisitCount,
			ReplyCount:         article.ReplyCount,
			ArticleTitle:       article.ArticleTitle,
			ArticlePreviewText: article.ArticlePreviewText,
			ArticleContent:     article.ArticleContent,
			Top:                article.Top,
			Category:           category,
			TagList:            tagList,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"data":       articlesRes,
	})

	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
	}).Info("Get articles success")
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
	lastID, err := models.AddArticle(article)
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
		"id":      lastID,
		"message": "Add article success",
	})

	// 打印相关日志
	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
		"title":      articleVals.ArticleTitle,
		"create_at":  time.Now().Format("2006-01-02 15:04:05"),
	}).Info("Add article success")

	// 把新增博文同步更新到 redis 里, 每次 mysql 更新成功都把所有 articles 同步到 redis
	articles := new([]models.Article)
	for count := 0; count < 3; count++ {
		articles, err = models.GetAllArticle()
		if err == nil && len(*articles) != 0 {
			break
		}

		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Get all article failed")
	}

	for _, value := range *articles {
		ma, _ := json.Marshal(value)

		// 向链表头push，因为都是倒序查找
		_ = models.RedisClient.LPush("article", ma)
	}

	// ma, _ := json.Marshal(models.Article{
	// 	ID:                 uint(lastID),
	// 	CreateAt:           article.CreateAt,
	// 	UpdateAt:           article.UpdateAt,
	// 	VisitCount:         article.VisitCount,
	// 	ReplyCount:         article.ReplyCount,
	// 	ArticleTitle:       article.ArticleTitle,
	// 	ArticlePreviewText: article.ArticlePreviewText,
	// 	ArticleContent:     article.ArticleContent,
	// 	Top:                false,
	// 	Category:           article.Category,
	// 	TagList:            article.TagList,
	// })
	// err = models.RedisClient.RPush("articles", ma).Err()
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"errorMsg": err,
	// 		"article":  article,
	// 	}).Info("Sync article to redis failed")

	// 	return
	// }

	log.WithFields(log.Fields{}).Info("Sync article to redis success")
}

// 根据数据库的category id返回相应的分类
func categoryForRes(id uint, categories *[]models.Category) models.Category {
	var category models.Category
	log.Info(*categories)

	for _, value := range *categories {
		log.Info(id)
		log.Info(value)
		if value.ID == id {
			category = value

			break
		}
	}

	return category
}

// 根据数据的 tag_list 返回相应的标签组
func tagForRes(tagStr string, tags *[]models.Tag) ([]models.Tag, error) {
	var tagSlice []models.Tag
	tagStrSlice := strings.Split(tagStr, "_")
	for _, tagID := range tagStrSlice {
		t, err := strconv.ParseUint(tagID, 10, 64)
		if err != nil {
			return nil, err
		}
		for _, tag := range *tags {
			ut := uint(t)
			if tag.ID == ut {
				tagSlice = append(tagSlice, tag)
			}
		}

	}

	return tagSlice, nil
}
