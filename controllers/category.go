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
	Category string
}

// 获取全部分类名的数据结构
// type Categories []string

// AddCategory 新增一个分类名
func AddCategory(c *gin.Context) {
	var categoryVals CategoryForm

	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 检查是否存在 category 字段
	err = c.ShouldBindWith(&categoryVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message":    err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Add category failed")

		return
	}

	// 向数据库添加分类名
	err = models.AddCategory(categoryVals.Category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message":    http.StatusText(http.StatusInternalServerError),
			"statusCode": http.StatusInternalServerError,
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store category to database failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Add category failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

	// 打印成功增加分类名的日志
	log.WithFields(log.Fields{
		"message":    "Add category success",
		"category":   categoryVals.Category,
		"statusCode": http.StatusOK,
	}).Info("Add category success")
}

// 获取全部分类名
func GetAllCategoryHandler(c *gin.Context) {
	categories, err := models.GetAllCategory()
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "query failed",
		})
		c.AbortWithError(401, errors.New("query failed"))

		return
	}

	log.Println(categories)

	c.JSON(200, gin.H{
		"categories": categories,
	})
}

// 删除某个分类名
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
