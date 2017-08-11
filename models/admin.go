package models

import (
	//"fmt"
	//"crypto/md5"
	"golang.org/x/crypto/bcrypt"
	"fmt"
	"log"
)

//设置后台登录账号密码长度范围
const (
	AdminIDLengthMin = 5
	AdminIDLengthMax = 20
	AdminPasswordLengthMin = 5
	AdminPasswordLengthMax = 20
)

//sql查询语句
const (
	qCreateInitAdmin = "INSERT INTO admin (admin_id, password, admin_name) VALUES (?, ?, ?)"
	qAdminByAdminID = "SELECT id, admin_id, password, admin_name, email, image FROM admin WHERE admin_id=?"
	qAdminByID = "SELECT id, admin_id, password, admin_name, email, image FROM admin WHERE id = ?"
	qAll = "SELECT * FROM admin"
)

//后台登录账号密码表单结构
type LoginAdminForm struct {
	AdminID		string `json:"admin_id"`
	Password	string 	`json:"password"`
}

type Admin struct {
	ID			uint		`db:"id" json:"id"`
	AdminID		string		`db:"admin_id" json:"admin_id"`
	Password	string		`db:"password" json:"password"`
	AdminName 	string		`db:"admin_name" json:"admin_name"`
	Email		string		`db:"email" json:"email"`
	Image		string		`db:"image" json:"image"`
	Token		string      `db:"_" json:"token"`
}

//admin请求响应
type AdminResponse struct {
	Admin  *Admin  `json:"admin"`
}

//初始化管理员账号
func (sdb *ServerDB) CreateInitAdmin() {
	hp, err := hashPassword("admin")
	if err != nil {
		log.Fatal(err)
	}
	a := Admin{
		AdminID: "admin",
		Password: string(hp),
		AdminName: "admin",
	}

	stmt, err := sdb.DB.Prepare(qCreateInitAdmin)
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(a.AdminID, a.Password, a.AdminName)
	if err != nil {
		log.Fatal(err)
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lastID)
}

//根据管理员id从mysql查询
func (sdb *ServerDB) AdminByAdminID(admin_id string) (*Admin, error){
	var a Admin
	if err := sdb.DB.QueryRow(qAdminByAdminID, admin_id).Scan(&a.ID, &a.AdminID, &a.Password, &a.AdminName, &a.Email, &a.Image); err != nil {
		return &a, err
	}
	return &a, nil
}

//加密密码，数据库不存储明文密码
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//验证密码
func (a *Admin) ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
}