/*
* admin controller
*
* token 1天后过期。每次使用延长6个小时，最长使用时效为一周
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
	"golang.org/x/crypto/bcrypt"
)

// SigningAlgorithm token加密算法
var SigningAlgorithm = "HS256"

// secret key
var secretKey = []byte("adminblog")

// Timeout token持续时间, 设置为一周
var Timeout = time.Hour * 24 * 7

// var Timeout = time.Second * 3600

// AdminLoginForm 登录表单
type AdminLoginForm struct {
	AdminID  string `form:"admin_id" json:"admin_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AdminUpdatePasswordForm 更改密码表单
type AdminUpdatePasswordForm struct {
	OldPassword string `form:"old_password" json:"old_password" binding:"required"`
	NewPassword string `form:"new_password" json:"new_password" binding:"required"`
}

// AdminUpdateInfoForm 更改管理员信息表单
type AdminUpdateInfoForm struct {
	AdminID   string `form:"admin_id" json:"admin_id" binding:"required"`
	AdminName string `form:"admin_name" json:"admin_name" binding:"required"`
	Image     string `form:"image" json:"image" binding:"required"`
}

// AdminLoginHandler 管理员登录，没有限制同一个ip的错误登录次数，容易被爆破，以后加登录序列
func AdminLoginHandler(c *gin.Context) {
	var loginVals AdminLoginForm

	// 判断是否有密码字段
	err := c.ShouldBindWith(&loginVals, binding.JSON)
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Miss adminID or password",
		})
		c.AbortWithStatus(http.StatusBadRequest)

		// 记录ip地址
		ip := c.ClientIP()
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"loginIp":    ip,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 检验 adminID
	adminID, err := checkString(loginVals.AdminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect adminID format",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 检验password
	password, err := checkString(loginVals.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect password format",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 检验 adminID 是否规定长度
	ab := checkAdminIDLength(adminID)
	if !ab {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect adminID length",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 检验 password 是否规定长度
	pb := checkAdminPasswordLength(password)
	if !pb {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect password length",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 从数据库查询密码
	admin, err := models.AdminByID()

	// 数据查询失败
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Admin information query failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin login failed")

		return
	}

	// 管理员ID不存在
	if admin.AdminID != adminID {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "AdminID not exist",
		})
		c.AbortWithStatus(http.StatusBadRequest)

		// 记录ip地址便于检测是否有人恶意爆破
		ip := c.ClientIP()
		log.WithFields(log.Fields{
			"errorMsg":   "AdminID not exist",
			"statusCode": http.StatusBadRequest,
			"ip":         ip,
		}).Info("Admin login failed")

		return
	}

	// 密码不匹配
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect admin password",
		})
		c.AbortWithStatus(http.StatusBadRequest)

		// 记录ip地址以便检测是否有人恶意爆破管理员账号
		ip := c.ClientIP()
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect admin password",
			"ip":         ip,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 生成token
	token := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	// 设置token过期时间
	expire := time.Now().Add(Timeout)
	claims["id"] = loginVals.AdminID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()

	// 生成token
	tokenString, err := token.SignedString(secretKey)

	// 生成token失败
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Generate token failed",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Generate token failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin login failed")

		return
	}

	// 把管理员 token 放入 redis里, token不同步进数据库
	err = models.RedisClient.Set("admin_token", tokenString, time.Hour*24).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Generate token failed",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store token to redis failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin login failed")

		return
	}

	// 生成token成功
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "login success",
		"token":      tokenString,
		"max_expire": expire.Format(time.RFC3339),
	})

	// 记录成功登录时间和 ip
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	ip := c.ClientIP()
	log.WithFields(log.Fields{
		"id":         loginVals.AdminID,
		"login_at":   currentTime,
		"ip":         ip,
		"statusCode": http.StatusOK,
	}).Info("Admin login success")

	// 存储登录成功的时间和 ip
	err = models.UpdateTimeIP(currentTime, ip)
	if err != nil {
		log.WithFields(log.Fields{
			"id":       loginVals.AdminID,
			"login_at": currentTime,
			"ip":       ip,
		}).Info("Store admin login time and ip failed")
	}
}

// GetAdminInfo 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	// token 认证
	admin, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 验证成功，返回管理员信息
	c.JSON(http.StatusOK, gin.H{
		"statusCode":    http.StatusOK,
		"admin_name":    admin.AdminName,
		"image":         admin.Image,
		"last_login_in": admin.LastLoginAt,
		"ip":            admin.IP,
	})

	// 打印日志
	log.WithFields(log.Fields{
		"admin_name":    admin.AdminName,
		"image":         admin.Image,
		"last_login_in": admin.LastLoginAt,
		"ip":            admin.IP,
		"statusCode":    http.StatusOK,
	}).Info("Get admin infomation success")
}

// AdminLogout 管理员退出
func AdminLogout(c *gin.Context) {
	// token 认证
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 把 redis 里的 admin_token 设置为空字符串
	err = models.RedisClient.Set("admin_token", "", 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Admin log out failed",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin login out failed")

		return
	}

	// 管理员成功退出
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Admin log out success",
	})

	log.WithFields(log.Fields{
		"message":    "Admin log out success",
		"statusCode": http.StatusOK,
	}).Info("Admin login out success")
}

// AdminUpdatePasswordhandler 更改管理员密码
func AdminUpdatePasswordhandler(c *gin.Context) {
	var adminUpdatePasswordVals AdminUpdatePasswordForm

	// token 认证
	admin, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 检查是否绑定了 password field
	err = c.ShouldBindWith(&adminUpdatePasswordVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "请填写新旧密码",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Miss password",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin change password failed")

		return
	}

	// 检查密码是否规定
	oldPassword, oerr := checkString(adminUpdatePasswordVals.OldPassword)
	newPassword, nerr := checkString(adminUpdatePasswordVals.NewPassword)
	if oerr != nil || nerr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "密码必须为数字字母组成",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect password format",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin change password failed")

		return
	}

	opb := checkAdminPasswordLength(oldPassword)
	nbp := checkAdminPasswordLength(newPassword)
	if !opb || !nbp {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "密码长度不正确",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect password length",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin change password failed")

		return
	}

	// 验证旧密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "旧密码验证不正确",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect password",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin change password failed")

		return
	}

	// 更改成新密码
	err = models.UpdateAdminPassword(newPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "密码修改失败",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin change password failed")

		return
	}

	err = models.RedisClient.Set("admin_token", "", 0).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "密码修改失败",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin change password failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "密码修改成功",
	})

	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
	}).Info("Admin change password success")
}

// AdminUpdateInfoHandler 更改管理员信息
func AdminUpdateInfoHandler(c *gin.Context) {
	var adminUpdateInfoVals AdminUpdateInfoForm

	// token 认证
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	// 检查是否绑定了 password field
	err = c.ShouldBindWith(&adminUpdateInfoVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Miss some field",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("Admin change information failed")

		return
	}

	// 把更新的管理员数据写入数据库
	err = models.UpdateAdminInfo(adminUpdateInfoVals.AdminID, adminUpdateInfoVals.AdminName, adminUpdateInfoVals.Image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "Change admin information failed",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Change admin information failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Change admin information success",
	})

	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
	}).Info("Change admin information success")
}

// AdminUploadImageHandler 上传管理员头像,不允许超过 4M
func AdminUploadImageHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Upload image failed",
		})
		log.WithFields(log.Fields{
			"message":    err,
			"statusCode": http.StatusBadRequest,
		}).Info("Upload image failed")

		return
	}

	err = c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Upload image failed",
		})
		log.WithFields(log.Fields{
			"message":    err,
			"statusCode": http.StatusBadRequest,
		}).Info("Upload image failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Upload image success",
	})

	log.WithFields(log.Fields{
		"statusCode": http.StatusOK,
	}).Info("Upload image success")
}

// 除去两边空格并检测字符串是否由数字和字母组成
func checkString(str string) (string, error) {
	// 两边除去换行符和空格
	s := strings.TrimSpace(str)

	// 判断s是否由数字和字母组成
	matched, err := regexp.MatchString("^[A-Za-z0-9]+$", s)

	if !matched || err != nil {
		return "", errors.New("String don't match regexp")
	}

	return s, nil
}

// 检测管理员 ID 是否规定长度
func checkAdminIDLength(adminID string) bool {
	if len(adminID) < models.AdminIDLengthMin || len(adminID) > models.AdminIDLengthMax {
		return false
	}

	return true
}

// 检测 password 是否规定长度
func checkAdminPasswordLength(password string) bool {
	if len(password) < models.AdminPasswordLengthMin || len(password) > models.AdminPasswordLengthMax {
		return false
	}

	return true
}
