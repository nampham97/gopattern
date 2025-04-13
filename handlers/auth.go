// File: handlers/auth.go
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("my_secret_key")

// Giả lập database người dùng với role
var userStore = map[string]struct {
	Password string
	Role     string
}{
	"admin": {Password: "password", Role: "admin"},
	"alice": {Password: "alice123", Role: "editor"},
	"bob":   {Password: "bob123", Role: "viewer"},
}

func generateTokens(username, role string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(2 * time.Minute).Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		"username": username,
		"type":     "refresh",
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims).SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims).SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	user, exists := userStore[creds.Username]
	if !exists || user.Password != creds.Password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := generateTokens(creds.Username, user.Role)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(req.RefreshToken, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil || !token.Valid {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["type"] != "refresh" {
		http.Error(w, "Invalid token type", http.StatusBadRequest)
		return
	}

	username, ok := claims["username"].(string)
	if !ok {
		http.Error(w, "Invalid token payload", http.StatusBadRequest)
		return
	}

	user, exists := userStore[username]
	if !exists {
		http.Error(w, "User no longer exists", http.StatusUnauthorized)
		return
	}

	newAccessToken, _, err := generateTokens(username, user.Role)
	if err != nil {
		http.Error(w, "Failed to generate new token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

// Route dùng thử để test quyền role
func AdminOnly(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(jwt.MapClaims)
	role := claims["role"].(string)
	if role != "admin" {
		http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome Admin!",
	})
}

// Route mẫu kiểm tra role sử dụng middleware
func EditorOnly(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Welcome Editor!",
	})
}
