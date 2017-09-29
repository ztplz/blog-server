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
var Timeout = time.Hour * 24 * 7

// AdminLogin  管理员登录表单
type AdminLogin struct {
	AdminID  string `form:"admin_id" json:"admin_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// AdminAuthMiddleware 后台token认证中间件
func AdminAuthMiddleware(c *gin.Context) (*models.Admin, error) {
	token, err := parseToken(c)

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

	return admin, nil
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
func parseToken(c *gin.Context) (*jwt.Token, error) {
	var token string
	var err error

	parts := strings.Split("header:Authorization", ":")
	token, err = jwtFromHeader(c, parts[1])
	if err != nil {
		return nil, err
	}

	// 查询 redis 里是否有这个token
	adminToken, err := models.RedisClient.Get("admin_token").Result()
	if err != nil {
		return nil, err
	}

	// 比对token
	if token != adminToken {
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
