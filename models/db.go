/*
* 数据库初始化连接
*
* author: ztplz
* email: mysticzt@gmail.com
* github: https://github.com/ztplz
* create-at: 2017.08.15
 */

package models

import (
	"database/sql"

	// _ "github.com/go-sql-driver/mysql"  因为放在models包，所以不在main函数里引入
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// DB  定义为sql.DB类型
var DB *sql.DB

// InitialDB 数据库初始化连接函数
func InitialDB() {
	// 连接数据库
	db, err := sql.Open("mysql", "root:123456@/blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		// log.Fatal(err)
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Database connect failed")
	}

	// 测试数据库是否正常练级
	if err := db.Ping(); err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Database ping failed")
	}

	DB = db
}
