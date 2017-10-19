package models

import (
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// User  定义User结构体
type User struct {
	ID          uint   `db:"id" json:"id"`
	UserID      string `db:"user_id" json:"user_id"`
	Password    string `db:"password" json:"password"`
	UserName    string `db:"user_name" json:"user_name"`
	Image       string `db:"image" json:"image"`
	CreateAt    string `db:"create_at" json:"create_at"`
	LastLoginAt string `db:"last_login_at" json:"last_login_at"`
	LoginCount  uint   `db:"login_count" json:"login_count"`
	IsBlacklist bool   `db:"is_blacklist" json:"is_blacklist"`
}

const (
	// UserIDLengthMin 用户ID 最少长度
	UserIDLengthMin = 5

	// UserIDLengthMax 用户ID 最大长度
	UserIDLengthMax = 15

	// UserNameLengthMin 用户名 最少长度
	UserNameLengthMin = 3

	// UserNameLengthMax 用户名 最大长度
	UserNameLengthMax = 21

	// UserPasswordLengthMin 用户密码 最少长度
	UserPasswordLengthMin = 5

	// UserPasswordLengthMax 用户密码 最大长度
	UserPasswordLengthMax = 20
)

const (
	qGetAllUser = "SELECT id, user_id, user_name, image, last_login_at, login_count, is_blacklist FROM user"
	qInsertUser = `INSERT INTO user 
						(user_id, user_name, password, image, create_at, last_login_at, login_count, is_blacklist)
						VALUES
						(?, ?, ?, ?, ?, ?, ?, ?)`
	qGetUserByUserID    = "SELECT id, user_id, user_name, password, image, create_at, last_login_at, login_count, is_blacklist FROM user WHERE user_id = ?"
	qUpdateUserID       = "UPDATE user SET user_id = ? WHERE user_id = ?"
	qUpdateUserName     = "UPDATE user SET user_name = ? WHERE user_name = ?"
	qUpdateUserPassword = "UPDATE user SET password = ? WHERE user_id = ?"
)

// GetAllUser 获取所有注册用户信息
func GetAllUser() (*[]User, error) {
	var u User
	var users []User

	// 根据管理员ID查询管理员信息
	rows, err := DB.Query(qGetAllUser)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Query all user failed")

		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		// 这里以后优化性能
		err = rows.Scan(
			&u.ID,
			&u.UserID,
			&u.UserName,
			&u.Image,
			&u.CreateAt,
			&u.LastLoginAt,
			&u.LoginCount,
			&u.IsBlacklist)
		users = append(users, u)
	}
	err = rows.Err()
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Rows scan failed")

		return nil, err
	}

	return &users, nil
}

// GetUserByUserID 根据 UserID 从数据库获取数据
func GetUserByUserID(userID string) (*User, error) {
	var user User

	row := DB.QueryRow(qGetUserByUserID, userID)
	err := row.Scan(&user.ID, &user.UserID, &user.UserName, &user.Password, &user.Image, &user.CreateAt, &user.LastLoginAt, &user.LoginCount, &user.IsBlacklist)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
			"userID":   userID,
		}).Info("Query user failed")

		return nil, err
	}

	return &user, nil
}

// UpdateUserID 更改用户ID
func UpdateUserID(oldUserID string, newUserID string) error {
	stmt, err := DB.Prepare(qUpdateUserID)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUserID, oldUserID)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql exec failed")

		return err
	}

	return nil
}

// UpdateUserName 更改用户名字
func UpdateUserName(oldUserName string, newUserName string) error {
	stmt, err := DB.Prepare(qUpdateUserName)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(newUserName, oldUserName)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql exec failed")

		return err
	}

	return nil
}

// UpdateUserPassword 更新用户密码
func UpdateUserPassword(userID string, password string) error {
	// 加密密码
	hp, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Encrypt user password failed")

		return err
	}

	stmt, err := DB.Prepare(qUpdateUserPassword)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(hp, userID)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Update user password failed")

		return err
	}

	return nil
}

// UserRegister 用户信息存储进数据库
func UserRegister(user *User) error {
	// 加密密码
	hp, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Encrypt user password failed")

		return err
	}

	stmt, err := DB.Prepare(qInsertUser)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Sql prepare failed")

		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UserID, user.UserName, string(hp), user.Image, user.CreateAt, user.LastLoginAt, user.LoginCount, user.IsBlacklist)
	if err != nil {
		log.WithFields(log.Fields{
			"errorMsg": err,
		}).Info("Insert user info to database failed")

		return err
	}

	return nil
}
