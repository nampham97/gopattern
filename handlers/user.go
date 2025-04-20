package handlers

import (
	"GoPattern/internal/logger"
	"GoPattern/internal/params"
	"GoPattern/internal/redisdb"
	"GoPattern/repository"
	"context"
	"encoding/json"
	"fmt"

	"net/http"
	"strconv"
	"strings"
	"time"

	"GoPattern/models"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type UserHandler struct {
	repo  *repository.UserRepository
	redis *redisdb.RedisClient
}

func NewUserHandler(repo *repository.UserRepository, redis *redisdb.RedisClient) *UserHandler {
	return &UserHandler{
		repo:  repo,
		redis: redis,
	}
}

// file: handlers/user.go
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, _ := strconv.Atoi(query.Get("page"))
	if page <= 0 {
		page = 1
	}
	limit, _ := strconv.Atoi(query.Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	filter := models.UserFilter{
		Search: query.Get("search"),
		Role:   query.Get("role"),
		Page:   page,
		Limit:  limit,
	}
	fmt.Println("filter:", filter)
	users, err := h.repo.GetUsers(filter)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	// Lấy ID từ URL và chuyển đổi sang int
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Tạo context với timeout 5 giây
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Gọi repository để lấy user
	user, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error fetching user: %v", err), http.StatusInternalServerError)
		return
	}

	// Trả về user dưới dạng JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetPaginatedUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	users, err := h.repo.GetPaginatedUsers(page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetUsers_filter(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := logger.FromContext(ctx)
	logger.Info("Fetching users")

	q := params.ParseQuery(r.URL.Query())

	users, err := h.repo.FindUsers(q)
	if err != nil {
		logger.Error("Failed to fetch", zap.Error(err))
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}
