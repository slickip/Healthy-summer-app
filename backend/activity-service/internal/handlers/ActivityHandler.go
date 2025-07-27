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
	UserID         uint   `json:"user_id"` // если используешь JWT, можно брать из токена
	ActivityTypeID uint   `json:"activity_type_id"`
	Duration       int    `json:"duration"`
	Intensity      string `json:"intensity"`
	Calories       int    `json:"calories"`
}

func (h *ActivityHandler) ActiveHandler(w http.ResponseWriter, r *http.Request) {
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

	// Теперь ждем activity_type_id вместо type
	var req struct {
		ActivityTypeID uint   `json:"activity_type_id"`
		Duration       int    `json:"duration"`
		Intensity      string `json:"intensity"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Достаем ActivityType из базы
	var at models.ActivityType
	if err := h.DB.First(&at, req.ActivityTypeID).Error; err != nil {
		http.Error(w, "Invalid activity type", http.StatusBadRequest)
		return
	}

	// Считаем калории
	totalCalories := req.Duration * at.CaloriesPerMinute

	activity := models.Activity{
		UserID:         userIDFromContext,
		ActivityTypeID: req.ActivityTypeID,
		Duration:       req.Duration,
		Intensity:      req.Intensity,
		Calories:       totalCalories,
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
		Preload("ActivityType").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&activities).Error; err != nil {
		http.Error(w, "Failed to fetch activities", http.StatusInternalServerError)
		return
	}

	type ActivityResponse struct {
		ID        uint   `json:"id"`
		Type      string `json:"type"`
		Duration  int    `json:"duration"`
		Intensity string `json:"intensity"`
		Calories  int    `json:"calories"`
	}

	var resp []ActivityResponse
	for _, a := range activities {
		resp = append(resp, ActivityResponse{
			ID:        a.ID,
			Type:      a.ActivityType.Name,
			Duration:  a.Duration,
			Intensity: a.Intensity,
			Calories:  a.Calories,
		})
	}

	json.NewEncoder(w).Encode(resp)
}
