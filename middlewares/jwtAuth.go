/*
* jwt token 验证模块
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package middlewares

import (
	"errors"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	// "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/ztplz/blog-server/models"
)

// SigningAlgorithm token加密算法
var SigningAlgorithm = "HS256"

// secret key
var secretKey = []byte("adminblog")

// Timeout token持续时间, 设置为一周
// var Timeout = time.Hour * 24 * 7
var Timeout = time.Second

// AdminLogin  管理员登录表单
type AdminLogin struct {
	AdminID  string `form:"admin_id" json:"admin_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AdminAuthMiddleware 后台token认证中间件
func AdminAuthMiddleware(c *gin.Context) (*models.Admin, error) {
	token, err := parseToken(c, "admin")

	// 如果解析 token 发生错误
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    err.Error(),
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusUnauthorized,
		}).Info("Admin auth failed")

		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"].(string)
	c.Set("JWT_PAYLOAD", claims)
	c.Set("AdminID", id)

	// 从数据取出管理员 ID
	admin, err := models.AdminByID()
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"message":    http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithStatus(http.StatusInternalServerError)
		log.WithFields(log.Fields{
			"errorMsg":   "admin info query failed",
			"statusCode": http.StatusInternalServerError,
		}).Info("Admin auth failed")

		return nil, err
	}

	if id != admin.AdminID {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    "You don't have permission to access",
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		log.WithFields(log.Fields{
			"errorMsg":   "Id don't match adminID",
			"statusCode": http.StatusUnauthorized,
		}).Info("Admin auth failed")

		return nil, errors.New("Incorrect token")
	}

	// redis延长6个小时 token 存储时间
	rtoken, err := models.RedisClient.Get("admin_token").Result()
	if err != nil {
		return nil, err
	}

	err = models.RedisClient.Set("admin_token", rtoken, time.Hour*6).Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Admin token store to redis failed")
	}

	return admin, nil
}

// UserAuthMiddleware 后台token认证中间件
func UserAuthMiddleware(c *gin.Context, userID string) error {
	_, err := parseToken(c, userID)

	// 如果解析 token 发生错误
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"message":    err.Error(),
		})
		c.AbortWithStatus(http.StatusUnauthorized)
		log.WithFields(log.Fields{
			"errorMsg":   err,
			"statusCode": http.StatusUnauthorized,
		}).Info("Admin auth failed")

		return err
	}

	// rtoken, err := models.RedisClient.Get(userID + "_token").Result()
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"statusCode": http.StatusUnauthorized,
	// 		"message":    "请重新登录",
	// 	})
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	log.WithFields(log.Fields{
	// 		"errorMsg":   err,
	// 		"statusCode": http.StatusUnauthorized,
	// 	}).Info("Admin auth failed")

	// 	return err
	// }

	// if *token != rtoken {
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"statusCode": http.StatusUnauthorized,
	// 		"message":    "你没有权限查看",
	// 	})
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	log.WithFields(log.Fields{
	// 		"errorMsg":   err,
	// 		"statusCode": http.StatusUnauthorized,
	// 	}).Info("Admin auth failed")

	// 	return err
	// }

	// claims := token.Claims.(jwt.MapClaims)

	// id := claims["id"].(string)
	// c.Set("JWT_PAYLOAD", claims)
	// c.Set("id", id)

	// // 从数据取出用户 ID
	// user, err := models.AdminByID()
	// if err != nil {
	// 	c.Header("WWW-Authenticate", "JWT realm=gin jwt")
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"statusCode": http.StatusInternalServerError,
	// 		"message":    http.StatusText(http.StatusInternalServerError),
	// 	})
	// 	c.AbortWithStatus(http.StatusInternalServerError)
	// 	log.WithFields(log.Fields{
	// 		"errorMsg":   "user info query failed",
	// 		"statusCode": http.StatusInternalServerError,
	// 	}).Info("User auth failed")

	// 	return nil, err
	// }

	// if id != user.UserID {
	// 	c.Header("WWW-Authenticate", "JWT realm=gin jwt")
	// 	c.JSON(http.StatusUnauthorized, gin.H{
	// 		"statusCode": http.StatusUnauthorized,
	// 		"message":    "You don't have permission to access",
	// 	})
	// 	c.AbortWithStatus(http.StatusUnauthorized)
	// 	log.WithFields(log.Fields{
	// 		"errorMsg":   "Id don't match adminID",
	// 		"statusCode": http.StatusUnauthorized,
	// 	}).Info("Admin auth failed")

	// 	return nil, errors.New("Incorrect token")
	// }

	// redis延长6个小时 token 存储时间

	// err = models.RedisClient.GetSet(userID+"_token", rtoken, time.Hour*6).Err()
	// if err != nil {
	// 	log.WithFields(log.Fields{
	// 		"errorMsg": err,
	// 	}).Info("Admin token store to redis failed")
	// }

	return nil
}

// ExtractClaims 提取 JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {
	_, exists := c.Get("JWT_PAYLOAD")
	if !exists {
		emptyClaims := make(jwt.MapClaims)

		return emptyClaims
	}

	jwtClaims, _ := c.Get("JWT_PAYLOAD")

	return jwtClaims.(jwt.MapClaims)
}

// 解析token
func parseToken(c *gin.Context, key string) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split("header:Authorization", ":")
	token, err = jwtFromHeader(c, parts[1])
	if err != nil {
		return nil, err
	}

	// 查询 redis 里是否有这个token
	rToken, err := models.RedisClient.Get(key + "_token").Result()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Admin auth failed")

		return nil, errors.New("token is expired")
	}

	// 如果 redis 里token为空值
	if rToken == "" {
		return nil, errors.New("You don't login in or token expired")
	}

	// 比对token
	if token != rToken {
		return nil, errors.New("Token is not match")
	}

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != token.Method {
			err := errors.New("Invalid signing algorithm")

			return nil, err
		}

		return secretKey, nil
	})
}

// 从请求头提取 token
func jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)
	log.Info(authHeader)

	// 如果请求头 Authorization 部分为空
	if authHeader == "" {
		err := errors.New("Auth header empty")

		return "", err
	}

	// 要求使用 Bearer Token
	parts := strings.SplitN(authHeader, " ", 2)

	if !(len(parts) == 2 && parts[0] == "Bearer") {
		err := errors.New("Invalid auth header")

		return "", err
	}

	return parts[1], nil
}
