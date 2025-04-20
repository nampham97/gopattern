package utils

import (
	"testing"
)

func TestCheckPasswordHash(t *testing.T) {
	// Mật khẩu và hash từ userStore
	password := "password"
	hash := "$2a$14$sRiEFgCBsvbT3VHR6Yuq3.s9d9tZtqDu2Ayj.EHxblLh1Eazdm4DS"

	// Kiểm tra xem hàm có trả về true không
	if !CheckPasswordHash(password, hash) {
		t.Errorf("CheckPasswordHash failed for correct password")
	}

	// Kiểm tra với mật khẩu sai
	if CheckPasswordHash("wrong_password", hash) {
		t.Errorf("CheckPasswordHash returned true for incorrect password")
	}
}
