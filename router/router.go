package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ztplz/blog-server/controllers"
)

// just for test token auth
// func helloHandler(c *gin.Context) {
// 	claims := controllers.ExtractClaims(c)
// 	c.JSON(200, gin.H{
// 		"userID": claims["id"],
// 		"text":   "Hello World.",
// 	})
// }

// 初始化路由
func InitialRouter() {
	r := gin.Default()

	// 后台管理登录
	r.POST("/api/v1/admin/login", controllers.AdminLoginHandler)

	// 管理员权限
	admin := r.Group("/api/v1/admin")
	{
		admin.GET("", controllers.GetAdminInfo)
		admin.POST("", controllers.AdminLoginHandler)
	}

	// 分类名操作
	category := r.Group("/api/v1/category")
	{
		category.GET("", controllers.GetAllCategoryHandler)
		category.POST("", controllers.AddCategory)
		category.DELETE("", controllers.DeleteCategoryHandler)
		category.PUT("/:name", controllers.UpdateCategoryHandler)
	}

	// 标签操作
	tag := r.Group("/api/v1/tag")
	{
		tag.GET("", controllers.GetAllTagHandler)
		tag.POST("", controllers.AddTagHandler)
	}

	r.Run(":8080")
}
