package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword băm mật khẩu với bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // cost 14 là khá an toàn
	return string(bytes), err
}

// CheckPasswordHash so sánh mật khẩu nhập vào với hash đã lưu
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
