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

	// 设置单次上传文件最大限制, 允许超过 8M
	r.MaxMultipartMemory = 8 << 20

	r.Use(middlewares.CORSMiddleware())

	// 统计访问人数
	r.Use(middlewares.CountVisitorMiddleware())

	// 管理员权限
	admin := r.Group("/api/v1/admin")
	{
		// 获取管理员信息，需鉴定token
		admin.GET("", controllers.GetAdminInfo)

		// 后台登录
		admin.POST("", controllers.AdminLoginHandler)

		// 更改管理员密码
		admin.PUT("/password", controllers.AdminUpdatePasswordhandler)

		// 更改管理员信息
		admin.PUT("", controllers.AdminUpdateInfoHandler)

		// 上传管理员头像
		admin.PUT("/image", controllers.AdminUploadImageHandler)

		// 管理员退出后台
		admin.DELETE("", controllers.AdminLogout)
	}

	// 博文操作
	article := r.Group("/api/v1/articles")
	{
		//获取所有文章列表， limit(每次返回列表数) page(页数)
		article.GET("", controllers.GetAllArticlesHandler)

		// 增加博文
		article.POST("", controllers.AddArticleHandler)
	}

	// 分类名操作
	category := r.Group("/api/v1/categories")
	{
		// 获取全部分类名 （article）true 查询每个分类的文章数
		category.GET("", controllers.GetAllCategoryHandler)

		// 增加分类名
		category.POST("", controllers.AddCategoryHandler)

		// category.DELETE("", controllers.DeleteCategoryHandler)
		category.PUT("/:name", controllers.UpdateCategoryHandler)
	}

	// 标签操作
	tag := r.Group("/api/v1/tags")
	{
		// 获取所有标签
		tag.GET("", controllers.GetAllTagHandler)

		// 增加标签
		tag.POST("", controllers.AddTagHandler)

		// 修改某个标签
		tag.PUT("/:id", controllers.UpdateTagHandler)
	}

	// 用户操作
	user := r.Group(".api/v1/user")
	{
		// 获取所有用户信息
		user.GET("", controllers.GetAllUser)

		// 获取某个用户的信息
		// user.GET("/:userID", controllers.GetUserByUserID)

		// 用户注册
		user.POST("", controllers.RegisterUser)
	}

	// 访客操作
	visitor := r.Group("api/v1/visitor")
	{
		visitor.GET("/count", controllers.GetVisitCount)
	}

	// 监听8080端口
	r.Run(":8080")
}
