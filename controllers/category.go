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

// CategoryForm 增加分类名表单结构
type CategoryForm struct {
	Category string `form:"category" json:"category" binding:"required"`
}

// 获取全部分类名的数据结构
// type Categories []string

// AddCategoryHandler 新增一个分类名
func AddCategoryHandler(c *gin.Context) {
	var categoryVals CategoryForm

	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 检查是否存在 category 字段
	err = c.ShouldBindWith(&categoryVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Miss category",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Add category failed")

		return
	}

	// 向数据库添加分类名
	lastID, err := models.AddCategory(categoryVals.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store category to database failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Add category failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "success",
	})

	// 打印成功增加分类名的日志
	log.WithFields(log.Fields{
		"message":    "Add category success",
		"category":   categoryVals.Category,
		"statusCode": http.StatusOK,
	}).Info("Add category success")

	// 向redis 写入 此次数据
	err := models.RedisClient.HSet("categories", string(lastID), models.Category{ID: uint(lastID), Category: catgoryVals}, 0).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"id":       lastID,
			"category": categoryVals,
		}).Info("Store category to redis failed")
	}
}

// GetAllCategoryHandler 获取全部分类名或者分类名包含的博文
func GetAllCategoryHandler(c *gin.Context) {
	// 从redis 里获取所有 category 的值
	categories, err := models.RedisClient.HVals("category").Result()

	// 如果从 redis 查询失败或者不存在 就从数据库读取全部分类名返回给用户, 并同步缺失的数据到 redis 里
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Get all categories from redis failed")

		// 从数据库获取全部分类名
		categories, err := models.GetAllCategory()
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    http.StatusText(http.StatusInternalServerError),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"message":    "Query categories failed",
				"statusCode": http.StatusInternalServerError,
			}).Info("Get all categories failed")

			return
		}

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"categories": categories,
		})

		// 同步到 redis 里
		for _, category := range categories {
			err := models.RedisClient.hset("categories", string(category.ID), category).Result()
			if err != nil {
				log.WithFields(log.Fields{
					"errorMsg": err,
					"category": category,
				}).Info("Sync category to redis failed")
			}
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"categories": categories,
	})

	log.WithFields(log.Fields{
		"categories": categories,
		"statusCode": http.StatusOK,
	}).Info("Get all categories success")
}

// DeleteCategoryHandler 删除某个分类名
func DeleteCategoryHandler(c *gin.Context) {
	var categoryVals CategoryForm

	_, err := middlewares.AdminAuthMiddleware(c)

	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		c.AbortWithError(401, errors.New("auth failed"))

		return
	}

	err = c.ShouldBindWith(&categoryVals, binding.JSON)
	if err != nil {
		log.Println(err)
		c.JSON(400, gin.H{
			"message": "Miss category field",
		})
		c.AbortWithError(400, errors.New("Miss category field"))

		return
	}

	err = models.DeleteCategory(categoryVals.Category)
	log.Println(err)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "failed",
		})
		c.AbortWithError(400, errors.New("failed"))

		return
	}

	c.JSON(200, gin.H{
		"message": "success",
	})
}

// 修改某个分类名
func UpdateCategoryHandler(c *gin.Context) {
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
	category := c.Query("name")
	// 替换原来的分类名
	key := c.Query("category")

	err = models.UpdateCategory(category, key)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "failed",
		})
		c.AbortWithError(401, errors.New("update failed"))

		return
	}

	c.JSON(200, gin.H{
		"message": "update success",
	})

	// err = c.ShouldBindWith(&categoryVals, binding.JSON)
	// if err != nil {
	// 	log.Println(err)
	// 	c.JSON(400, gin.H{
	// 		"message": "Miss category field",
	// 	})
	// 	c.AbortWithError(400, errors.New("Miss category field"))

	// 	return
	// }

	// err = models.DeleteCategory(categoryVals.Category)
	// log.Println(err)
	// if err != nil {
	// 	c.JSON(500, gin.H{
	// 		"message": "failed",
	// 	})
	// 	c.AbortWithError(400, errors.New("failed"))

	// 	return
	// }

	// c.JSON(200, gin.H{
	// 	"message": "success",
	// })
}
