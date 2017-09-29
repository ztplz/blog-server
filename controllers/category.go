/*
* 分类名处理
*
* 分类名以 hash 的格式存储在 redis 里，
* 并且以在 mysql 中的主键作为存在 redis 里的 field 名，
* 以次来映射 redis 里分类名和 mysql 里的分类名的关系
*
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package controllers

import (
	"encoding/json"
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

	// 检查是否有重复分类
	b, err := checkRepeatCategory("categories", categoryVals.Category)
	if err != nil {
		// 直接返回错误码500
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Add category failed")

		return
	}

	// 如果已经存在一样的分类名
	if b {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Category already exist",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Category already exist",
			"category":   categoryVals.Category,
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
		"lastID":     lastID,
		"category":   categoryVals.Category,
		"statusCode": http.StatusOK,
	}).Info("Add category success")

	// 向redis 写入 此次数据
	category, _ := json.Marshal(models.Category{ID: uint(lastID), Category: categoryVals.Category})
	// log.Info(string(category))
	// log.WithFields(log.Fields{
	// 	"stringcategory": string(category),
	// 	"category":       category,
	// }).Info("asdfasgdads")
	// strcategory := new(models.Category)
	// _ = json.Unmarshal(category, &strcategory)
	// log.WithFields(log.Fields{
	// 	"category": *strcategory,
	// }).Info("asdfasgdads")
	err = models.RedisClient.HSet("categories", string(lastID), category).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"id":       lastID,
			"errorMsg": err,
			"category": categoryVals,
		}).Info("Store category to redis failed")
	}
}

// GetAllCategoryHandler 获取全部分类名或者分类名包含的博文
func GetAllCategoryHandler(c *gin.Context) {
	// 从redis 里获取所有 category 的值
	categories, err := models.RedisClient.HVals("categories").Result()

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

		// 返回数据给客户端
		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"categories": categories,
		})

		// 同步到 redis 里
		for _, category := range categories {
			err := models.RedisClient.HSet("categories", string(category.ID), category).Err()
			if err != nil {
				log.WithFields(log.Fields{
					"errorMsg": err,
					"category": category,
				}).Info("Sync category to redis failed")
			}
		}

		return
	}

	// gin 有个很大的问题, JSON方法会把body里的再次序列化, 所以这里把数据从 redis 里取出来后又被 JSON方法 序列化
	cts := new([]models.Category)
	for _, value := range categories {
		ct := new(models.Category)
		err := json.Unmarshal([]byte(value), &ct)
		if err != nil {
			c.JSON(500, gin.H{
				"statusCode": http.StatusInternalServerError,
				"message":    http.StatusText(http.StatusInternalServerError),
			})
			c.AbortWithStatus(http.StatusInternalServerError)
			log.WithFields(log.Fields{
				"errorMsg":   err,
				"statusCode": http.StatusInternalServerError,
			}).Info("Unmasrshal category failed")

			return
		}

		*cts = append(*cts, *ct)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"categories": *cts,
	})

	log.WithFields(log.Fields{
		"categories": *cts,
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
}

// 就检查是否重复提交已存在分类名
func checkRepeatCategory(key string, field string) (bool, error) {
	// 首先开始从 redis 里查询
	categories, err := models.RedisClient.HVals(key).Result()
	if err != nil {
		// 从数据库里查询
		categories, err := models.GetAllCategory()
		if err != nil {
			return false, err
		}

		for _, category := range categories {
			if category.Category == field {
				return true, nil
			}
		}

		return false, nil
	}

	ct := models.Category{}
	for _, value := range categories {
		err := json.Unmarshal([]byte(value), &ct)
		if err != nil {
			return false, err
		}

		// 判断是否和以前的分类名相同
		if ct.Category == field {
			return true, nil
		}
	}

	return false, nil
}
