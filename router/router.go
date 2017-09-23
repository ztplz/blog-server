/*
* blog 路由
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ztplz/blog-server/controllers"
	"github.com/ztplz/blog-server/middlewares"
)

// InitialRouter  初始化路由
func InitialRouter() {
	r := gin.Default()

	r.Use(middlewares.CORSMiddleware())

	// 管理员权限
	admin := r.Group("/api/v1/admin")
	{
		// 获取管理员信息，需鉴定token
		admin.GET("", controllers.GetAdminInfo)

		// 后台登录
		admin.POST("", controllers.AdminLoginHandler)
	}

	// 博文操作
	article := r.Group("/api/v1/article")
	{
		article.POST("", controllers.AddArticleHandler)
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

	// 监听8080端口
	r.Run(":8080")
}
