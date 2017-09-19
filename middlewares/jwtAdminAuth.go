package middlewares

import (
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/ztplz/blog-server/models"
)

// token加密算法
var SigningAlgorithm = "HS256"

// secret key
var secretKey = []byte("adminblog")

// token持续时间, 设置为一周
var Timeout = time.Hour * 24 * 7

type AdminLogin struct {
	AdminID  string `form:"admin_id" json:"admin_id" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// 后台token认证中间件
func AdminAuthMiddleware(c *gin.Context) (*models.Admin, error) {
	token, err := parseToken(c)

	// 如果解析 token 发生错误
	if err != nil {
		return nil, err
		// unauthorized(c, http.StatusUnauthorized, err.Error())
		// return
	}

	claims := token.Claims.(jwt.MapClaims)

	id := claims["id"].(string)
	c.Set("JWT_PAYLOAD", claims)
	c.Set("AdminID", id)

	// 从数据取出管理员 ID
	admin, err := models.AdminByID(1)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithError(http.StatusInternalServerError, errors.New("Query AdminID Faild"))

		return nil, err
	}

	if id != admin.AdminID {
		unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return nil, errors.New("AdminID Do Not Match")
	}

	return admin, nil
}

// 提取 JWT claims
func ExtractClaims(c *gin.Context) jwt.MapClaims {

	if _, exists := c.Get("JWT_PAYLOAD"); !exists {
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

	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod(SigningAlgorithm) != token.Method {
			cause := errors.New("Invalid Signing Algorithm")
			err := errors.WithMessage(cause, "Auth Failed")
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
		cause := errors.New("Auth Header Empty")
		err := errors.WithMessage(cause, "Auth Failed")
		return "", err
	}

	// 要求使用 Bearer Token
	parts := strings.SplitN(authHeader, " ", 2)

	if !(len(parts) == 2 && parts[0] == "Bearer") {
		cause := errors.New("Invalid Auth Header")
		err := errors.WithMessage(cause, "Auth Failed")
		return "", err
	}

	return parts[1], nil
}

func unauthorized(c *gin.Context, code int, message string) {
	c.Header("WWW-Authenticate", "JWT realm=gin jwt")
	c.JSON(code, gin.H{
		"message": message,
	})
	c.AbortWithError(code, errors.New(message))

	return
}
