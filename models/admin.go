package models

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
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
	qCreateInitAdmin = "INSERT INTO admin (admin_id, password, admin_name, image, last_login_at, ip) VALUES (?, ?, ?, ?, ?, ?)"
	qAdminByAdminID  = "SELECT id, admin_id, password, admin_name, image, last_login_at FROM admin WHERE admin_id=?"
	qAdminByID       = "SELECT id, admin_id, password, admin_name, image, last_login_at, ip FROM admin WHERE id=?"
	// qAdminByID = "SELECT id, admin_id, password, admin_name, email, image FROM admin WHERE id = ?"
	// qAll = "SELECT * FROM admin"
	qUpdateLastLoginAt = "UPDATE admin SET last_login_at = ?, ip = ? WHERE id = 1"
)

// LoginAdminForm 后台登录账号密码表单结构
// type LoginAdminForm struct {
// 	AdminID  string `json:"admin_id"`
// 	Password string `json:"password"`
// }

// Admin  定义Admin结构体
type Admin struct {
	ID          uint   `db:"id" json:"id"`
	AdminID     string `db:"admin_id" json:"admin_id"`
	Password    string `db:"password" json:"password"`
	AdminName   string `db:"admin_name" json:"admin_name"`
	Image       string `db:"image" json:"image"`
	LastLoginAt string `db:"last_login_at" json:"last_login_at"`
	IP          string `db:"ip" json:"ip"`
}

// InitialAdmin 初始化管理员账号
func InitialAdmin() {
	// 加密初始管理员密码
	hp, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("encrypt initial password failed")
	}

	// 存储初始管理员账号信息
	stmt, err := DB.Prepare(qCreateInitAdmin)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Sql prepare failed")
	}
	defer stmt.Close()

	_, err = stmt.Exec("admin", string(hp), "admin", "", "", "")
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Fatal("Insert initial password to database failed")
	}

	log.Info("Create initial password success")
}

// AdminByID 从数据库根据 ID 查询管理员信息
func AdminByID() (*Admin, error) {
	var a Admin

	// 根据管理员ID查询管理员信息
	err := DB.QueryRow(qAdminByID, 1).Scan(&a.ID, &a.AdminID, &a.Password, &a.AdminName, &a.Image, &a.LastLoginAt, &a.IP)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Query admin infomation failed")

		return &a, err
	}

	return &a, nil
}

// UpdateTimeIP 存储管理员登录时间
func UpdateTimeIP(lastLoginAt string, ip string) error {
	// sql预处理
	stmt, err := DB.Prepare(qUpdateLastLoginAt)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg":      err,
			"last_login_at": lastLoginAt,
			"ip":            ip,
		}).Info("Sql prepare admin last login time and ip failed")

		return err
	}

	// sql执行
	_, err = stmt.Exec(lastLoginAt, ip)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg":      err,
			"last_login_at": lastLoginAt,
			"ip":            ip,
		}).Info("Insert admin last login time and ip to database failed")

		return err
	}

	return nil
}
