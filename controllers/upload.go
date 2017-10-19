package controllers

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const accessKeyID string = "************************"

const accessKeySecret string = "************************"

const host string = "************************"

const expireTime int64 = 60

const uploadDir string = "/"

const (
	base64Table = "123QRSTUabcdVWXYZHijKLAWDCABDstEFGuvwxyzGHIJklmnopqr234560178912"
)

var coder = base64.NewEncoding(base64Table)

// ConfigStruct 配置
type ConfigStruct struct {
	Expiration string     `json:"expiration"`
	Conditions [][]string `json:"conditions"`
}

// PolicyToken 安全token
type PolicyToken struct {
	AccessKeyID string `json:"accessid"`
	Host        string `json:"host"`
	Expire      int64  `json:"expire"`
	Signature   string `json:"signature"`
	Policy      string `json:"policy"`
	Directory   string `json:"dir"`
}

// GetAccesskey 获取临时上传 accesskey
func GetAccesskey(c *gin.Context) {
	response := getPolicyToken()

	c.Header("Access-Control-Allow-Methods", "POST")
	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{
		"response": response,
	})
}

func base64Encode(src []byte) []byte {
	return []byte(coder.EncodeToString(src))
}

func getGmtISO8601(expireEnd int64) string {
	var tokenExpire = time.Unix(expireEnd, 0).Format("2006-01-02T15:04:05Z")

	return tokenExpire
}

// 获取临时token
func getPolicyToken() PolicyToken {
	var config ConfigStruct
	var condition []string
	var policyToken PolicyToken

	now := time.Now().Unix()
	expireEnd := now + expireTime

	tokenExpire := getGmtISO8601(expireEnd)

	//
	config.Expiration = tokenExpire

	condition = append(condition, "starts-with")
	condition = append(condition, "$key")
	condition = append(condition, uploadDir)
	config.Conditions = append(config.Conditions, condition)

	// 生成签名
	result, err := json.Marshal(config)
	debyte := base64.StdEncoding.EncodeToString(result)
	h := hmac.New(func() hash.Hash { return sha1.New() }, []byte(accessKeySecret))
	io.WriteString(h, debyte)
	signedStr := base64.StdEncoding.EncodeToString(h.Sum(nil))

	policyToken.AccessKeyID = accessKeyID
	policyToken.Host = host
	policyToken.Expire = expireEnd
	policyToken.Signature = string(signedStr)
	policyToken.Directory = uploadDir
	policyToken.Policy = string(debyte)
	// response, err := json.Marshal(policyToken)
	if err != nil {
		fmt.Println("json err:", err)
	}

	// return string(response)
	return policyToken
}
