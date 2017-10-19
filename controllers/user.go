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
)

// UserRegisterForm 用户注册表单结构
type UserRegisterForm struct {
	UserID   string `form:"user_id" json:"user_id" binding:"required"`
	UserName string `form:"user_name" json:"user_name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// GetAllUser 查询所有用户信息
func GetAllUser(c *gin.Context) {
	// 管理员鉴权
	_, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		return
	}

	users, err := models.GetAllUser()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Get all user info failed")

		return
	}

	// 查询成功
	c.JSON(http.StatusInternalServerError, gin.H{
		"statusCode": http.StatusInternalServerError,
		"message":    "success",
		"users":      *users,
	})

	log.WithFields(log.Fields{
		"message":    "Get all user info success",
		"statusCode": http.StatusOK,
	}).Info("Get all user info success")
}

// GetUserByUserID 根据 UserID 获取用户信息
func GetUserByUserID(c *gin.Context) {
	userID := c.Param("userID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户ID不能为空",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "用户ID为空",
			"statusCode": http.StatusBadRequest,
		}).Info("Get User failed")

		return
	}

	// 鉴权
	err := middlewares.UserAuthMiddleware(c, userID)
	if err != nil {
		return
	}

	user, err := models.GetUserByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "数据获取失败，请重试",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("Get User failed")

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "数据获取成功",
		"user":       *user,
	})
}

// RegisterUser 用户注册
func RegisterUser(c *gin.Context) {
	var userVals UserRegisterForm

	// 判断是否有必须字段
	err := c.ShouldBindWith(&userVals, binding.JSON)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "请完整填写相应的注册信息",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 除去两边空格
	userID, ierr := checkUserString(userVals.UserID)
	userName := strings.TrimSpace(userVals.UserName)
	// userName, nerr := checkUserString(userVals.UserName)
	password, perr := checkUserString(userVals.Password)
	if ierr != nil || perr != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "请不要在提交的注册信息中包含空格",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"ierr":       ierr,
			"perr":       perr,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 检查用户 ID 和用户密码是否由 数字 字母组成
	ib := checkStringCharNum(userID)
	pb := checkStringCharNum(password)
	if !ib || !pb {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户ID和密码必须有字母、数字组成",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 检查用户 ID 长度
	b := checkUserIDLength(userID)
	if !b {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户ID不符合规定长度",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 检查用户名长度
	b = checkUserNameLength(userName)
	if !b {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户名不符合规定长度",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 检查用户密码长度
	b = checkUserPasswordLength(userName)
	if !b {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户密码不符合规定长度",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusBadRequest,
		}).Info("User register failed")

		return
	}

	// 用户信息校验成功，生成token并把用户信息存入数据库，token放入redis
	token := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	// 设置token过期时间
	expire := time.Now().Add(Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = time.Now().Unix()

	// 生成token
	tokenString, err := token.SignedString(secretKey)

	// 生成token失败
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Generate token failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("User register failed")

		return
	}

	// 把用户 token 放入 redis里, token不同步进数据库
	key := userID + "_token"
	err = models.RedisClient.Set(key, tokenString, time.Hour*24).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store token to redis failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("User register failed")

		return
	}

	// 存储进数据
	user := &models.User{
		ID:          0,
		UserID:      userID,
		Password:    password,
		UserName:    userName,
		Image:       "",
		CreateAt:    time.Now().Format("2006-01-02 15:04:05"),
		LastLoginAt: time.Now().Format("2006-01-02 15:04:05"),
		LoginCount:  1,
		IsBlacklist: false,
	}

	err = models.UserRegister(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "Store user info to databse failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("User register failed")

		return
	}

	// 注册成功
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "注册成功",
		"userID":     userID,
		"token":      tokenString,
		"max_expire": expire.Format(time.RFC3339),
	})

	log.WithFields(log.Fields{
		"userID":     userID,
		"userName":   userName,
		"statusCode": http.StatusOK,
	}).Info("User register success")
}

// UpdateUserID 修改用户ID
func UpdateUserID(c *gin.Context) {
	oldUserID := c.Param("userID")
	newUserID := c.Query("new_userID")

	if oldUserID == "" || newUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户ID不能为空",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "用户ID为空",
			"statusCode": http.StatusBadRequest,
		}).Info("User userID update failed")

		return
	}

	// 鉴权
	err := middlewares.UserAuthMiddleware(c, oldUserID)
	if err != nil {
		return
	}

	err = models.UpdateUserID(oldUserID, newUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "更改失败，请重新尝试",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("User userID update failed")

		return
	}

	// 更新redis里token key
	token, _ := models.RedisClient.Get(oldUserID + "_token").Result()
	_ = models.RedisClient.Set(newUserID+"_token", token, 6).Err()

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "更改成功",
		"user_id":    newUserID,
	})
}

// UpdateUserName 更改用户名
func UpdateUserName(c *gin.Context) {
	oldUserName := c.Query("old_userName")
	newUserName := c.Query("new_userName")
	userID := c.Param("userID")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "你没权权限操作",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		log.WithFields(log.Fields{
			"errorMsg":   "userID 为空",
			"statusCode": http.StatusUnauthorized,
		}).Info("User userName update failed")

		return
	}

	if oldUserName == "" || newUserName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": http.StatusBadRequest,
			"message":    "用户名字不能为空",
		})
		c.AbortWithStatus(http.StatusBadRequest)
		log.WithFields(log.Fields{
			"errorMsg":   "用户名字为空",
			"statusCode": http.StatusBadRequest,
		}).Info("User update failed")

		return
	}

	// 鉴权
	err := middlewares.UserAuthMiddleware(c, userID)
	if err != nil {
		return
	}

	err = models.UpdateUserName(oldUserName, newUserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "更改失败，请重新尝试",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusInternalServerError,
		}).Info("User userName update failed")

		return
	}

	// 更新redis里的token key
	// token, _ := models.RedisClient.Get(userID + "_token").Result()
	// _ = models.RedisClient.Set(newUserID+"_token", token, 6).Err()

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "更改成功",
		"user_name":  newUserName,
	})
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(c *gin.Context) {
	userID := c.Param("userID")
	password := c.Query("password")

	err := models.UpdateUserPassword(userID, password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    "更改失败，请重新尝试",
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"userID":     userID,
			"statusCode": http.StatusInternalServerError,
		}).Info("Update user password failed")

		return
	}

	// 清除redis里原来的token
	_ = models.RedisClient.Del(userID + "_token")

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "更改成功",
	})
}

// 除去两边空格并检测输入里面是否有空格
func checkUserString(str string) (string, error) {
	// 两边除去换行符和空格
	s := strings.TrimSpace(str)

	// 判断字符串内是否有空格
	b := strings.Contains(s, " ")
	if b {
		return "", errors.New("String has space")
	}

	// 判断s是否由数字和字母组成
	matched, err := regexp.MatchString("^[A-Za-z0-9]+$", s)

	if !matched || err != nil {
		return "", errors.New("String don't match regexp")
	}

	return s, nil
}

// 检测用户id 和 密码是否有 数字 字母构成
func checkStringCharNum(str string) bool {
	matched, err := regexp.MatchString("^[A-Za-z0-9]+$", str)
	if !matched || err != nil {
		return false
	}

	return true
}

// 检测用户 ID 是否规定长度
func checkUserIDLength(userID string) bool {
	if len(userID) < models.UserIDLengthMin || len(userID) > models.UserIDLengthMax {
		return false
	}

	return true
}

// 检测用户名是否规定长度
func checkUserNameLength(userName string) bool {
	if len(userName) < models.UserNameLengthMin || len(userName) > models.UserNameLengthMax {
		return false
	}

	return true
}

// 检测 password 是否规定长度
func checkUserPasswordLength(password string) bool {
	if len(password) < models.UserPasswordLengthMin || len(password) > models.UserPasswordLengthMax {
		return false
	}

	return true
}
