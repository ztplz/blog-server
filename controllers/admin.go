package controllers

import (
	"github.com/ztplz/blog-server/models"
	"github.com/gin-gonic/gin"
	//"gin/json"

	"database/sql"
	"fmt"
	//"log"
)

type AdminController struct {

}

//处理后台登录请求
func (ctrl *AdminController) Login(c *gin.Context) {
	//从提交的form里提取账号和密码
	loginAdmin := &models.LoginAdminForm{AdminID: c.PostForm("admin_id"), Password: c.PostForm("password")}

	//检查账号密码是否为空值
	if loginAdmin.AdminID == "" || loginAdmin.Password == "" {
		fmt.Println("id or password nil")
		c.JSON(406, gin.H{
			"status": "账号或密码为空",
		})
		c.Abort()
		return
	}

	//检查账号密码是否过短
	if len(loginAdmin.AdminID) < models.AdminIDLengthMin || len(loginAdmin.Password) < models.AdminPasswordLengthMin {
		c.JSON(406, gin.H{
			"status": "账号或密码太短",
		})
		c.Abort()
		return
	}

	//检查账号密码是否过长
	if len(loginAdmin.AdminID) > models.AdminIDLengthMax || len(loginAdmin.Password) > models.AdminPasswordLengthMax {
		c.JSON(406, gin.H{
			"status": "账号或密码太长",
		})
		c.Abort()
		return
	}

	//和数据库里存储的密文对比验证
	a, err := models.SDB.AdminByAdminID(loginAdmin.AdminID)
	//数据库未匹配到后台管理员信息
	if err == sql.ErrNoRows {
		fmt.Println(err)
		c.JSON(406, gin.H{
			"status": "未匹配到管理员信息",
		})
		c.Abort()
		return
	}
	//其他错误
	if err != nil {
		fmt.Println(err)
		c.JSON(406, gin.H{
			"status": "数据库查询错误",
		})
		c.Abort()
		return
	}
	//如果结果为空值
	if a == nil {
		fmt.Println(err, "login: SDB.AdminByAdminID() a == nil")
		c.JSON(406, gin.H{
			"status": "未查询成功",
		})
		c.Abort()
		return
	}
	//账号密码验证错误
	if err := a.ValidatePassword(loginAdmin.Password); err != nil {
		fmt.Println(err)
		c.JSON(406, gin.H{
			"status": "账号或密码错误",
		})
		c.Abort()
		return
	}

	//登录成功，设置JWT TOKEN
	token, err := models.AdminNewToken(a)
	if err != nil {
		fmt.Println(err, "login: models.NewToken()")
	}
	a.Token = token

	//响应
	c.JSON(200, a)
}