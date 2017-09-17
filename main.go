package main

import (
	//"log"
	//"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ztplz/blog-server/controllers"
	"github.com/ztplz/blog-server/models"

	//"fmt"
)
func main() {
	db := models.ConnectDB()
	models.SDB.CreateInitAdmin()
	defer db.Close()

	//初始化路由
	router := gin.Default()

	//后台api
	v1_admin := router.Group("api/v1/admin")
	{
			admin := new(controllers.AdminController)

			//v1_admin.GET("", admin.GetAdminInfo)
			v1_admin.POST("", admin.Login)
	}

	router.Run(":8080")
}
