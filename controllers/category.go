package controllers

import (
	"errors"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
)

// 增加分类名表单结构
type CategoryForm struct {
	Category string
}

// 获取全部分类名的数据结构
// type Categories []string

// 新增一个分类名
func AddCategory(c *gin.Context) {
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

	err = models.AddCategory(categoryVals.Category)
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
	category := c.Param("name")
	// 替换原来的分类名
	key := c.Query("category")

	err = models.UpdateCategory(category, key)
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{
			"message": "failed",
		})
		c.AbortWithError(401, errors.New("update failed"))
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
