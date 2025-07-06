package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/models"
	"gorm.io/gorm"
)

type StepsHandler struct {
	DB *gorm.DB
}

type CreateStepsRequest struct {
	UserID    uint      `json:"user_id"`
	StepCount int       `json:"step_count"`
	Date      time.Time `json:"type:date"`
}

func (h *StepsHandler) StepHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateSteps(w, r)
	case http.MethodGet:
		h.ListSteps(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *StepsHandler) CreateSteps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		StepsCount int    `json:"steps_count"`
		Date       string `json:"date"` // "2025-07-06"
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	steps := models.Steps{
		UserID:    userID,
		StepCount: req.StepsCount,
		Date:      parsedDate,
	}

	if err := h.DB.Create(&steps).Error; err != nil {
		http.Error(w, "Failed to save steps", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(steps)
}

func (h *StepsHandler) ListSteps(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var steps []models.Steps
	if err := h.DB.
		Where("user_id = ?", userID).
		Order("date desc").
		Find(&steps).Error; err != nil {
		http.Error(w, "Failed to fetch steps", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(steps)
}
