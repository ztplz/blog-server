package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// 密码加密
func hashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// 密码验证
func ValidatePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(Password), []byte(password))
}