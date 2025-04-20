// File: handlers/auth.go
package handlers

import (
	base "GoPattern/internal/shared"
	"GoPattern/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	*base.BaseHandler
}

func NewAuthHandler(baseHandler *base.BaseHandler) *AuthHandler {
	return &AuthHandler{BaseHandler: baseHandler}
}

var jwtKey = []byte("my_secret_key")

// Giả lập database người dùng với role
var userStore = map[string]struct {
	Password string
	Role     string
}{
	"admin": {Password: "$2a$14$sRiEFgCBsvbT3VHR6Yuq3.s9d9tZtqDu2Ayj.EHxblLh1Eazdm4DS", Role: "admin"},  // "password"
	"alice": {Password: "$2a$14$1B5ePQTEc8aIcdgPavSC4Ox9X4CQK/HO7txa/X5hwCI8A.ENgrw.y", Role: "editor"}, // "alice123"
}

func generateTokens(username, role string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(10 * time.Minute).Unix(),
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

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	user, exists := userStore[creds.Username]
	if !exists || !utils.CheckPasswordHash(creds.Password, user.Password) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	accessToken, refreshToken, err := generateTokens(creds.Username, user.Role)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	h.RedisClient.SetRd(r.Context(), creds.Username, refreshToken, 10*time.Minute)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	// Parse token
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

	username := claims["username"].(string)

	// Kiểm tra token trong Redis
	storedToken, err := h.RedisClient.GetRd(r.Context(), username)
	if err != nil || storedToken != req.RefreshToken {
		http.Error(w, "Token mismatch", http.StatusUnauthorized)
		return
	}

	user, exists := userStore[username]
	if !exists {
		http.Error(w, "User no longer exists", http.StatusUnauthorized)
		return
	}

	newAccessToken, newRefreshToken, err := generateTokens(username, user.Role)
	if err != nil {
		http.Error(w, "Token generation failed", http.StatusInternalServerError)
		return
	}

	// Cập nhật token mới vào Redis
	err = h.RedisClient.SetRd(r.Context(), username, newAccessToken, 10*time.Minute)
	if err != nil {
		http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	err := h.RedisClient.DeleteRd(r.Context(), username)
	if err != nil {
		http.Error(w, "Logout failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out",
	})
}

// Route dùng thử để test quyền role
func AdminOnly(w http.ResponseWriter, r *http.Request) {
	role := r.Context().Value("role")
	if role == nil || role.(string) != "admin" {
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
