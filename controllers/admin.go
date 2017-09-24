package controllers

import (
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

// AdminLoginForm 登录表单
type AdminLoginForm struct {
	AdminID  string `form:"admin_id" json:"admin_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
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

	// 管理员账号，密码两边除去换行符和空格
	adminID := strings.TrimSpace(loginVals.AdminID)
	password := strings.TrimSpace(loginVals.Password)

	// 判断管理员账号，密码是否由数字和字母组成
	a, _ := regexp.MatchString("^[A-Za-z0-9]+$", adminID)
	b, _ := regexp.MatchString("^[A-Za-z0-9]+$", password)
	if !a || !b {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect format",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect format",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 判断管理员账户是否符合规定长度
	if len(adminID) < models.AdminIDLengthMin || len(adminID) > models.AdminIDLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect adminID Length",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect adminID Length",
			"statusCode": http.StatusBadRequest,
		}).Info("Admin login failed")

		return
	}

	// 判断密码是否符合规定长度
	if len(password) < models.AdminPasswordLengthMin || len(password) > models.AdminPasswordLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "Incorrect password length",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "Incorrect password length",
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

	// 生成token成功
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "login success",
		"token":      tokenString,
		"expire":     expire.Format(time.RFC3339),
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
	c.JSON(200, gin.H{
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
