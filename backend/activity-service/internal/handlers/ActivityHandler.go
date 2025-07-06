package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/activity-service/internal/models"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	DB *gorm.DB
}

// Request payload
type CreateActivityRequest struct {
	UserID    uint   `json:"user_id"` // если используешь JWT, можно брать из токена
	Type      string `json:"type"`
	Duration  int    `json:"duration"`
	Intensity string `json:"intensity"`
	Calories  int    `json:"calories"`
	Location  string `json:"location"`
}

func (h *ActivityHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateActivity(w, r)
	case http.MethodGet:
		h.ListActivities(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ActivityHandler) CreateActivity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userIDFromContext := userIDValue.(uint)

	var req CreateActivityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	activity := models.Activity{
		UserID:    userIDFromContext,
		Type:      req.Type,
		Duration:  req.Duration,
		Intensity: req.Intensity,
		Calories:  req.Calories,
	}

	if err := h.DB.Create(&activity).Error; err != nil {
		http.Error(w, "Failed to create activity", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(activity)
}

func (h *ActivityHandler) ListActivities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var activities []models.Activity
	if err := h.DB.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&activities).Error; err != nil {
		http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(activities)
}
