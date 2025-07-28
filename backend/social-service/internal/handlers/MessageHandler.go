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

type MessageHandler struct {
	DB *gorm.DB
}

// SendMessage godoc
// @Summary Send a message to another user
// @Description Send a private message to specified user
// @Tags messaging
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body models.Messages true "Message data"
// @Success 201 {object} models.Messages
// @Failure 400 {string} string "Invalid input data"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/messages [post]
func (h *ChallengeHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	senderID := userIDValue.(uint)

	var req struct {
		ReceiverID uint   `json:"receiver_id"`
		Content    string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.ReceiverID == 0 || req.Content == "" {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	message := models.Messages{
		SenderID:   senderID,
		RecieverID: req.ReceiverID,
		Content:    req.Content,
		SentAt:     time.Now(),
	}

	if err := h.DB.Create(&message).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(message)
}

// GetMessages godoc
// @Summary Get conversation history
// @Description Retrieve message history between authenticated user and specified friend
// @Tags messaging
// @Produce json
// @Security ApiKeyAuth
// @Param friend_id query int true "Friend ID to get conversation with"
// @Success 200 {array} models.Messages
// @Failure 400 {string} string "Missing or invalid friend_id"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/messages [get]
func (h *ChallengeHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	friendIDStr := r.URL.Query().Get("friend_id")
	friendID, err := strconv.Atoi(friendIDStr)
	if friendIDStr == "" || err != nil {
		http.Error(w, "Missing or invalid friend_id", http.StatusBadRequest)
		return
	}

	var messages []models.Messages
	if err := h.DB.Where(`
		(sender_id = ? AND reciever_id = ?) OR (sender_id = ? AND reciever_id = ?)`,
		userID, friendID, friendID, userID).Order("sent_at ASC").Find(&messages).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
