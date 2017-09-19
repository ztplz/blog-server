package controllers

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/ztplz/blog-server/middlewares"
	"github.com/ztplz/blog-server/models"
	"golang.org/x/crypto/bcrypt"
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

// 管理员登录，没有限制同一个ip的错误登录次数，容易被爆破，以后加登录序列
func AdminLoginHandler(c *gin.Context) {
	var loginVals AdminLogin

	ct := c.ContentType()
	log.Println("ct: " + ct)

	// 判断是否有密码字段
	if c.ShouldBindWith(&loginVals, binding.JSON) != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(400, gin.H{
			"message": "Miss AdminID Or Password",
		})
		c.AbortWithError(400, errors.New("Miss AdminID Or Password"))

		return
	}

	log.Println(loginVals.AdminID)
	log.Println(loginVals.Password)

	// 管理员账号，密码两边除去换行符和空格
	adminID := strings.TrimSpace(loginVals.AdminID)
	password := strings.TrimSpace(loginVals.Password)

	// 判断管理员账号，密码是否由数字和字母组成
	a, _ := regexp.MatchString("^[A-Za-z0-9]+$", adminID)
	b, _ := regexp.MatchString("^[A-Za-z0-9]+$", password)
	if !a || !b {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect Format",
		})
		c.AbortWithError(http.StatusBadRequest, errors.New("Incorrect Format"))

		return
	}

	// 判断管理员账户是否符合规定长度
	if len(adminID) < models.AdminIDLengthMin || len(adminID) > models.AdminIDLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect AdminID Length",
		})
		c.AbortWithError(http.StatusBadRequest, errors.New("Incorrect AdminID Length"))

		return
	}

	// 判断密码是否符合规定长度
	if len(password) < models.AdminPasswordLengthMin || len(password) > models.AdminPasswordLengthMax {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect Password Length",
		})
		c.AbortWithError(http.StatusBadRequest, errors.New("Incorrect Password Length"))

		return
	}

	// 从数据库查询密码
	admin, err := models.AdminByID(1)

	// 数据查询失败
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": http.StatusText(http.StatusInternalServerError),
		})
		c.AbortWithError(http.StatusInternalServerError, errors.New("Admin Data Query Fail"))

		return
	}

	// 管理员ID不存在
	if admin.AdminID != adminID {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "AdminID Not Exist",
		})
		c.AbortWithError(http.StatusInternalServerError, errors.New("AdminID Not Exist"))

		return
	}

	// 密码不匹配
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Incorrect password",
		})
		c.AbortWithError(http.StatusBadRequest, errors.New("Incorrect password"))

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
	log.Println(err)

	// 生成token失败
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "request token failed",
		})
		c.AbortWithError(http.StatusInternalServerError, errors.New("Create JWT Token faild"))

		return
	}

	// 生成token成功
	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"token":   tokenString,
		"expire":  expire.Format(time.RFC3339),
	})

	// 记录成功登录时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("管理员 %s 登录成功    %s", loginVals.AdminID, currentTime)
	err = models.UpdateAdminLoginTime(currentTime)
	if err != nil {
		log.Printf("存储管理员登录时间 %s 失败", currentTime)
	}
}

// 获取管理员信息
func GetAdminInfo(c *gin.Context) {
	// token 认证
	admin, err := middlewares.AdminAuthMiddleware(c)
	if err != nil {
		c.Header("WWW-Authenticate", "JWT realm=gin jwt")
		c.JSON(401, gin.H{
			"message": err.Error(),
		})
		c.AbortWithError(401, errors.New("auth failed"))

		return
	}

	c.JSON(200, gin.H{
		"admin_name":    admin.AdminName,
		"image":         admin.Image,
		"last_login_in": admin.LastLoginAt,
	})

}

// 管理员退出登录
func AdminLogOut(c *gin.Context) {

}

// 后台token认证中间件
func AdminAuthMiddleware(c *gin.Context) {
	token, err := parseToken(c)

	// 如果解析 token 发生错误
	if err != nil {
		unauthorized(c, http.StatusUnauthorized, err.Error())
		return
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
	}

	if id != admin.AdminID {
		unauthorized(c, http.StatusForbidden, "You don't have permission to access.")
		return
	}

	c.Next()
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
			return nil, errors.New("Invalid Signing Algorithm")
		}

		return secretKey, nil
	})
}

// 从请求头提取 token
func jwtFromHeader(c *gin.Context, key string) (string, error) {
	authHeader := c.Request.Header.Get(key)

	// 如果请求头 Authorization 部分为空
	if authHeader == "" {
		return "", errors.New("Auth Header Empty")
	}

	// 要求使用 Bearer Token
	parts := strings.SplitN(authHeader, " ", 2)
	log.Println(parts[0])
	log.Println(parts[1])

	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("Invalid Auth Header")
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
