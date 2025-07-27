package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/slickip/Healthy-summer-app/backend/nutrition-service/internal/models"
	"gorm.io/gorm"
)

type FoodHandler struct {
	DB *gorm.DB
}

func (h *FoodHandler) FoodHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateFood(w, r)
	case http.MethodGet:
		h.ListFood(w, r)
	case http.MethodDelete:
		h.DeleteFood(w, r)
	case http.MethodPut:
		h.UpdateFood(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FoodHandler) CreateFood(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name             string  `json:"name"`
		CalloriesPer100g float32 `json:"callories_per_100g"`
		Proteins         float32 `json:"proteins"`
		Fats             float32 `json:"fats"`
		Carbs            float32 `json:"carbs"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Логируем входящие данные
	fmt.Printf("Received food data: Name=%s, Calories=%.2f, Proteins=%.2f, Fats=%.2f, Carbs=%.2f\n",
		req.Name, req.CalloriesPer100g, req.Proteins, req.Fats, req.Carbs)

	food := models.Foods{
		Name:             req.Name,
		CalloriesPer100g: req.CalloriesPer100g,
		Proteins:         req.Proteins,
		Fats:             req.Fats,
		Carbs:            req.Carbs,
	}

	if err := h.DB.Create(&food).Error; err != nil {
		fmt.Printf("Error creating food: %v\n", err)
		http.Error(w, "Failed to create new food", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created food: ID=%d, Name=%s, Calories=%.2f, Proteins=%.2f, Fats=%.2f, Carbs=%.2f\n",
		food.ID, food.Name, food.CalloriesPer100g, food.Proteins, food.Fats, food.Carbs)

	// Формируем ответ с нужными полями
	type FoodResponse struct {
		ID               uint    `json:"id"`
		Name             string  `json:"name"`
		CalloriesPer100g float32 `json:"callories_per_100g"`
		Proteins         float32 `json:"proteins"`
		Fats             float32 `json:"fats"`
		Carbs            float32 `json:"carbs"`
	}

	response := FoodResponse{
		ID:               food.ID,
		Name:             food.Name,
		CalloriesPer100g: food.CalloriesPer100g,
		Proteins:         food.Proteins,
		Fats:             food.Fats,
		Carbs:            food.Carbs,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *FoodHandler) ListFood(w http.ResponseWriter, r *http.Request) {
	var foods []models.Foods
	if err := h.DB.Find(&foods).Error; err != nil {
		http.Error(w, "Failed to retrieve foods", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Retrieved %d foods from database:\n", len(foods))
	for i, food := range foods {
		fmt.Printf("  Food[%d]: ID=%d, Name=%s, Calories=%.2f, Proteins=%.2f, Fats=%.2f, Carbs=%.2f\n",
			i, food.ID, food.Name, food.CalloriesPer100g, food.Proteins, food.Fats, food.Carbs)
	}

	// Формируем ответ с нужными полями
	type FoodResponse struct {
		ID               uint    `json:"id"`
		Name             string  `json:"name"`
		CalloriesPer100g float32 `json:"callories_per_100g"`
		Proteins         float32 `json:"proteins"`
		Fats             float32 `json:"fats"`
		Carbs            float32 `json:"carbs"`
	}

	var resp []FoodResponse
	for _, food := range foods {
		resp = append(resp, FoodResponse{
			ID:               food.ID,
			Name:             food.Name,
			CalloriesPer100g: food.CalloriesPer100g,
			Proteins:         food.Proteins,
			Fats:             food.Fats,
			Carbs:            food.Carbs,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *FoodHandler) DeleteFood(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}

	if err := h.DB.Delete(&models.Foods{}, id).Error; err != nil {
		http.Error(w, "Failed to delete food", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *FoodHandler) UpdateFood(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
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
		Name             *string  `json:"name,omitempty"`
		CalloriesPer100g *float32 `json:"callories_per_100g,omitempty"`
		Proteins         *float32 `json:"proteins,omitempty"`
		Fats             *float32 `json:"fats,omitempty"`
		Carbs            *float32 `json:"carbs,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var food models.Foods
	if err := h.DB.First(&food, id).Error; err != nil {
		http.Error(w, "Food not found", http.StatusNotFound)
		return
	}

	if req.Name != nil {
		food.Name = *req.Name
	}
	if req.CalloriesPer100g != nil {
		food.CalloriesPer100g = *req.CalloriesPer100g
	}
	if req.Proteins != nil {
		food.Proteins = *req.Proteins
	}
	if req.Fats != nil {
		food.Fats = *req.Fats
	}
	if req.Carbs != nil {
		food.Carbs = *req.Carbs
	}

	if err := h.DB.Save(&food).Error; err != nil {
		http.Error(w, "Failed to update food", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(food)
}
