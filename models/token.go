package models

import (
	// "github.com/gin-gonic/gin"
	// "fmt"
	"time"

    "github.com/dgrijalva/jwt-go"
)

//开发阶段用固定secret-key， 正式上线用flag获取
var tokenSecret = []byte("golang")

//设置Token为字符类型
type Token string

//TokenClaims
type AdminTokenClaims struct {
	Admin  	*Admin
	jwt.StandardClaims
}


//JWT过期时间
const jwtExpiryDuration = time.Hour * 24 * 7

func AdminNewToken(a *Admin) (string, error) {
	claims := AdminTokenClaims{
		a,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtExpiryDuration).Unix(),
			Issuer: "mysticzt-blog",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(tokenSecret)
	if err != nil {
		return "", err
	}
	return ss, nil
}

// func AdminValidateToken(tokenString string) (*AdminTokenClaims, error) {
// 	token, err := jwt.ParseWithClaims(tokenString, &AdminTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		_, ok := token.Method.(*jwt.SigningMethodHMAC)
// 		if !ok {
// 			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return tokenSecret, nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}
// 	claims, ok := token.Claims.(*AdminTokenClaims)
// 	if !ok || !token.Valid {
// 		return nil, fmt.Errorf("Token not valid")
// 	}
// 	return claims, nil
// }

