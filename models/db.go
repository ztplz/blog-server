package models

import (
	//"fmt"
	"log"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// 连接数据库
func InitialDB() {
	db, err := sql.Open("mysql", "root:123456@/blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	// 测试数据库是否正常练级
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	DB = db
}