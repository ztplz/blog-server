package models

import (
	//"fmt"
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"log"
)

type ServerDB struct {
	DB		*sql.DB
}

var SDB   *ServerDB

//连接数据库，上线后用flag命令行获取密码连接
func ConnectDB() *sql.DB{
	//连接数据库
	db, err := sql.Open("mysql", "root:123456@/blog?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}

	//测试数据库是否正常连接
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	SDB = &ServerDB{db }

	return db
}
