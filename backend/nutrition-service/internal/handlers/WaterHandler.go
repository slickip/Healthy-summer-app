package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/models"
	"gorm.io/gorm"
)

type WaterHandler struct {
	DB *gorm.DB
}

func (h *WaterHandler) WaterHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateWaterLog(w, r)
	case http.MethodGet:
		h.ListWaterLog(w, r)
	case http.MethodDelete:
		h.DeleteWaterLog(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *WaterHandler) CreateWaterLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	// читаем тело как []byte один раз
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// логируем сырое тело
	fmt.Println("🔵 Raw JSON body:", string(bodyBytes))

	// декодируем в структуру
	var req struct {
		VolumeML int    `json:"volume_ml"`
		LoggedAt string `json:"logged_at"`
	}
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// ⏰ парсим дату
	loggedAt, err := time.Parse(time.RFC3339, req.LoggedAt)
	if err != nil {
		http.Error(w, "Invalid logged_at format. Use RFC3339", http.StatusBadRequest)
		return
	}

	waterLog := models.WaterLogs{
		UserID:   userID,
		VolumeML: req.VolumeML,
		LoggedAt: loggedAt,
	}

	if err := h.DB.Create(&waterLog).Error; err != nil {
		http.Error(w, "Failed to create water log", http.StatusInternalServerError)
		return
	}
	type WaterResponse struct {
		ID       uint      `json:"id"`
		VolumeML int       `json:"volume_ml"`
		LoggedAt time.Time `json:"logged_at"`
	}

	response := WaterResponse{
		ID:       waterLog.ID,
		VolumeML: waterLog.VolumeML,
		LoggedAt: waterLog.LoggedAt,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *WaterHandler) ListWaterLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var water_log []models.WaterLogs
	if err := h.DB.Where("user_id = ?", userID).Order("logged_at DESC").Find(&water_log).Error; err != nil {
		http.Error(w, "Failed to retrieve water log", http.StatusInternalServerError)
		return
	}

	// Формируем ответ с нужными полями
	type WaterLogResponse struct {
		ID        uint      `json:"id"`
		UserID    uint      `json:"user_id"`
		VolumeML  int       `json:"volume_ml"`
		LoggedAt  time.Time `json:"logged_at"`
		CreatedAt time.Time `json:"created_at"`
	}

	var resp []WaterLogResponse
	for _, log := range water_log {
		resp = append(resp, WaterLogResponse{
			ID:        log.ID,
			UserID:    log.UserID,
			VolumeML:  log.VolumeML,
			LoggedAt:  log.LoggedAt,
			CreatedAt: log.CreatedAt,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *WaterHandler) DeleteWaterLog(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.WaterLogs{}).Error; err != nil {
		http.Error(w, "Failed to delete water log", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
