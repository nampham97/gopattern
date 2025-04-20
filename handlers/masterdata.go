package handlers

import (
	base "GoPattern/internal/shared"
	"GoPattern/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type MasterdataHandler struct {
	*base.BaseHandler
	repo *repository.ProvincesRepository
}

func NewMasterdataHandler(baseHandler *base.BaseHandler, repo *repository.ProvincesRepository) *MasterdataHandler {
	return &MasterdataHandler{
		BaseHandler: baseHandler,
		repo:        repo,
	}
}

func (h MasterdataHandler) GetCountries(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Kiểm tra cache trong Redis
	cached, err := h.RedisClient.GetRd(r.Context(), "masterdata_countries")
	if err == nil && cached != "" {
		msg := fmt.Sprintf("✅ From Redis:\n%s\n⏱️ Took: %v", cached, time.Since(start))
		w.Write([]byte(msg))
		return
	}

	// Giả lập truy vấn DB lâu
	time.Sleep(500 * time.Millisecond)
	dbData := "Vietnam,USA,Japan,France"

	// Lưu vào Redis cache (10 phút)
	_ = h.RedisClient.SetRd(r.Context(), "masterdata_countries", dbData, 10*time.Minute)

	msg := fmt.Sprintf("🛢️ From DB:\n%s\n⏱️ Took: %v", dbData, time.Since(start))
	w.Write([]byte(msg))
}

func (h *MasterdataHandler) GetProvinces(w http.ResponseWriter, r *http.Request) {
	// Kiểm tra cache trong Redis
	cached, err := h.RedisClient.GetRd(r.Context(), "masterdata_provinces")
	if err == nil && cached != "" {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		log.Println("data from redis:", cached)
		return
	}
	log.Println("data from db")
	// Tạo context với timeout 5 giây
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	provinces, err := h.repo.GetProvinces(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching provinces: %v", err), http.StatusInternalServerError)
		return
	}

	provincesJSON, _ := json.Marshal(provinces)
	_ = h.RedisClient.SetRd(r.Context(), "masterdata_provinces", string(provincesJSON), 10*time.Minute)

	if err := json.NewEncoder(w).Encode(provinces); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}
