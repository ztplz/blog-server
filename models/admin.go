package models

import (
	//"fmt"
	//"crypto/md5"
	"golang.org/x/crypto/bcrypt"
	// "fmt"
	"log"
)

//设置后台登录账号密码长度范围
const (
	AdminIDLengthMin       = 5
	AdminIDLengthMax       = 20
	AdminPasswordLengthMin = 5
	AdminPasswordLengthMax = 20
)

//sql查询语句
const (
	qCreateInitAdmin = "INSERT INTO admin (admin_id, password, admin_name, image, last_login_at) VALUES (?, ?, ?, ?, ?)"
	qAdminByAdminID  = "SELECT id, admin_id, password, admin_name, image, last_login_at FROM admin WHERE admin_id=?"
	qAdminByID       = "SELECT id, admin_id, password, admin_name, image, last_login_at FROM admin WHERE id=?"
	// qAdminByID = "SELECT id, admin_id, password, admin_name, email, image FROM admin WHERE id = ?"
	// qAll = "SELECT * FROM admin"
	qUpdateLastLoginAt = "UPDATE admin SET last_login_at=? WHERE id = 1"
)

//后台登录账号密码表单结构
type LoginAdminForm struct {
	AdminID  string `json:"admin_id"`
	Password string `json:"password"`
}

type Admin struct {
	ID        uint   `db:"id" json:"id"`
	AdminID   string `db:"admin_id" json:"admin_id"`
	Password  string `db:"password" json:"password"`
	AdminName string `db:"admin_name" json:"admin_name"`
	Image     string `db:"image" json:"image"`
	// Token		string      `db:"_" json:"token"`
	LastLoginAt string `db:"last_login_at" json:"last_login_at"`
}

// 初始化管理员账号
func InitialAdmin() {
	hp, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	log.Println(hp)
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := DB.Prepare(qCreateInitAdmin)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	log.Println("string hp: " + string(hp))

	_, err = stmt.Exec("admin", string(hp), "admin", "", "")
	// res, err := stmt.Exec(a.AdminID, a.Password, a.AdminName, a.Image, a.Token)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("create admin success")
}

// 从数据库根据 AdminID 查询管理员密码
func AdminByAdminID(admin_id string) (*Admin, error) {
	var a Admin

	if err := DB.QueryRow(qAdminByAdminID, admin_id).Scan(&a.ID, &a.AdminID, &a.Password, &a.AdminName, &a.Image, &a.LastLoginAt); err != nil {
		return &a, err
	}

	// log.Println(a.ID)
	log.Println("a.AdminID: " + a.AdminID)
	log.Println("a.Password: " + a.Password)
	log.Println("a.AdminName: " + a.AdminName)

	return &a, nil
}

// 从数据库跟 ID 查询管理员信息
func AdminByID(id uint) (*Admin, error) {
	var a Admin

	if err := DB.QueryRow(qAdminByID, id).Scan(&a.ID, &a.AdminID, &a.Password, &a.AdminName, &a.Image, &a.LastLoginAt); err != nil {
		return &a, err
	}

	return &a, nil
}

// 存储管理员登录时间
func UpdateAdminLoginTime(last_login_at string) error {
	// sql预处理
	stmt, err := DB.Prepare(qUpdateLastLoginAt)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(last_login_at)
	if err != nil {
		return err
	}

	log.Printf("update admin last login time success    s%", last_login_at)

	return nil
}
