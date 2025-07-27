package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/models"
	"gorm.io/gorm"
)

type MealHandler struct {
	DB *gorm.DB
}

func (h *MealHandler) Mealhandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateMeal(w, r)
	case http.MethodGet:
		h.ListMeal(w, r)
	case http.MethodDelete:
		h.DeleteMeal(w, r)
	case http.MethodPut:
		h.UpdateMeal(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *MealHandler) CreateMeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		MealTime    string `json:"meal_time"` // ожидаем ISO8601 строку
		Description string `json:"description"`
		Calories    int    `json:"calories"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	mealTime, err := time.Parse(time.RFC3339, req.MealTime)
	if err != nil {
		http.Error(w, "Invalid meal_time format. Use RFC3339", http.StatusBadRequest)
		return
	}

	meal := models.Meals{
		UserID:      userID,
		MealTime:    mealTime,
		Description: req.Description,
		Calories:    req.Calories,
	}

	if err := h.DB.Create(&meal).Error; err != nil {
		http.Error(w, "Failed to create meal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(meal)
}

func (h *MealHandler) ListMeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var meals []models.Meals
	if err := h.DB.Where("user_id = ?", userID).Order("meal_time DESC").Find(&meals).Error; err != nil {
		http.Error(w, "Failed to retrieve meals", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meals)
}

func (h *MealHandler) DeleteMeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	query := r.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id parameter", http.StatusBadRequest)
		return
	}

	if err := h.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Meals{}).Error; err != nil {
		http.Error(w, "Failed to delete meal", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *MealHandler) UpdateMeal(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	// Получаем id из query параметра
	query := r.URL.Query()
	idStr := query.Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	var req struct {
		MealTime    *string `json:"meal_time,omitempty"`
		Description *string `json:"description,omitempty"`
		Calories    *int    `json:"calories,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Находим запись
	var meal models.Meals
	if err := h.DB.First(&meal, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		http.Error(w, "Meal not found", http.StatusNotFound)
		return
	}

	// Обновляем поля, если они заданы
	if req.MealTime != nil {
		parsedTime, err := time.Parse(time.RFC3339, *req.MealTime)
		if err != nil {
			http.Error(w, "Invalid meal_time format", http.StatusBadRequest)
			return
		}
		meal.MealTime = parsedTime
	}
	if req.Description != nil {
		meal.Description = *req.Description
	}
	if req.Calories != nil {
		meal.Calories = *req.Calories
	}

	if err := h.DB.Save(&meal).Error; err != nil {
		http.Error(w, "Failed to update meal", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(meal)
}
