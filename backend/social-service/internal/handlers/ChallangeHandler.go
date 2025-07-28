package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/models"
	"gorm.io/gorm"
)

type ChallengeHandler struct {
	DB *gorm.DB
}

func (h *ChallengeHandler) ChallengeHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.CreateChallenge(w, r)
	case http.MethodGet:
		h.ListChallenges(w, r)
	case http.MethodPut:
		h.UpdateChallenge(w, r)
	case http.MethodDelete:
		h.DeleteChallenge(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateChallenge godoc
// @Summary Create new challenge
// @Description Create a new challenge by the authenticated user
// @Tags challenges
// @Accept  json
// @Produce  json
// @Param challenge body models.Challanges true "Challenge to create"
// @Success 201 {object} models.Challanges
// @Failure 400 {string} string "Invalid input"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/challenges [post]
func (h *ChallengeHandler) CreateChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var challenge models.Challanges
	if err := json.NewDecoder(r.Body).Decode(&challenge); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if challenge.Title == "" || challenge.Description == "" || challenge.GoalValue <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

	challenge.CreatorID = userID
	challenge.CreatedAt = time.Now()

	if err := h.DB.Create(&challenge).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(challenge)
}

// ListChallenges godoc
// @Summary List all challenges
// @Description Get a list of all available challenges
// @Tags challenges
// @Produce json
// @Success 200 {array} models.Challanges
// @Failure 500 {string} string "Internal server error"
// @Router /api/challenges [get]
func (h *ChallengeHandler) ListChallenges(w http.ResponseWriter, r *http.Request) {
	var challenges []models.Challanges
	if err := h.DB.Preload("ChallangesType").Find(&challenges).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenges)
}

// UpdateChallenge godoc
// @Summary Update a challenge
// @Description Update an existing challenge
// @Tags challenges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id query int true "Challenge ID"
// @Param challenge body models.Challanges true "Challenge update data"
// @Success 200 {object} models.Challanges
// @Failure 400 {string} string "Invalid ID or request body"
// @Failure 404 {string} string "Challenge not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/challenges [put]
func (h *ChallengeHandler) UpdateChallenge(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var challenge models.Challanges
	if err := h.DB.First(&challenge, id).Error; err != nil {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	var updateData models.Challanges
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.DB.Model(&challenge).Updates(updateData).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(challenge)
}

// DeleteChallenge godoc
// @Summary Delete a challenge
// @Description Delete an existing challenge
// @Tags challenges
// @Security ApiKeyAuth
// @Param id query int true "Challenge ID"
// @Success 204
// @Failure 400 {string} string "Invalid ID"
// @Failure 404 {string} string "Challenge not found"
// @Failure 500 {string} string "Internal server error"
// @Router /api/challenges [delete]
func (h *ChallengeHandler) DeleteChallenge(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result := h.DB.Delete(&models.Challanges{}, id)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	if result.RowsAffected == 0 {
		http.Error(w, "Challenge not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// JoinChallenge godoc
// @Summary Join a challenge
// @Description Join an existing challenge as a participant
// @Tags challenges
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.ChallangeParticipants true "Participant data"
// @Success 201 {object} models.ChallangeParticipants
// @Failure 400 {string} string "Invalid request body"
// @Failure 401 {string} string "Unauthorized"
// @Failure 409 {string} string "Already joined"
// @Failure 500 {string} string "Internal server error"
// @Router /api/challenges/join [post]
func (h *ChallengeHandler) JoinChallenge(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		ChallengeID uint `json:"challenge_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ChallengeID == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists := h.DB.Where("user_id = ? AND challange_type_id = ?", userID, req.ChallengeID).First(&models.ChallangeParticipants{}).Error == nil
	if exists {
		http.Error(w, "Already joined", http.StatusConflict)
		return
	}

	join := models.ChallangeParticipants{
		UserID:          int(userID),
		ChallangeTypeID: req.ChallengeID,
		Progress:        0,
		Status:          "began",
		JoinAt:          time.Now(),
	}

	if err := h.DB.Create(&join).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(join)
}

// MyChallenges godoc
// @Summary Get user's challenges
// @Description Get challenges created by the authenticated user
// @Tags challenges
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.Challanges
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/challenges/my [get]
func (h *ChallengeHandler) MyChallenges(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var challenges []models.Challanges
	if err := h.DB.Where("creator_id = ?", userID).Find(&challenges).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(challenges)
}

// ChallengeLeaderboard godoc
// @Summary Get leaderboard for a challenge
// @Description Returns ranked participants for a challenge
// @Tags challenges
// @Produce  json
// @Param id query int true "Challenge ID"
// @Success 200 {array} models.ChallangeParticipants
// @Router /api/challenges/leaderboard [get]
func (h *ChallengeHandler) ChallengeLeaderboard(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if idStr == "" || err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var leaderboard []models.ChallangeParticipants
	if err := h.DB.Where("challange_type_id = ?", id).
		Order("progress DESC").Find(&leaderboard).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(leaderboard)
}
