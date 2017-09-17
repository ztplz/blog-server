package router

import (
	"github.com/gin-gonic/gin"
	"github.com/ztplz/blog-server/controllers"
)

// just for test token auth
func helloHandler(c *gin.Context) {
	claims := controllers.ExtractClaims(c)
	c.JSON(200, gin.H{
		"userID": claims["id"],
		"text":   "Hello World.",
	})
}

// 初始化路由
func InitialRouter() {
	r := gin.Default()

	// 后台管理登录
	r.POST("/api/v1/admin/login", controllers.AdminLoginHandler)

	// 管理员权限
	authAdmin := r.Group("/api/v1/admin")
	authAdmin.Use(controllers.AdminAuthMiddleware)
	{
		authAdmin.GET("/hello", helloHandler)
	}

	r.Run(":8080")
}
