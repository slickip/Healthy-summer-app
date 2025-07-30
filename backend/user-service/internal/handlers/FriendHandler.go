package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/user-service/internal/models"
)

func (h *Handler) FriendHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/friends/request":
		if r.Method == http.MethodPost {
			h.SendFriendRequest(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	case "/api/friends/respond":
		if r.Method == http.MethodPost {
			h.HandleFriendRequest(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	case "/api/friends/requests":
		if r.Method == http.MethodGet {
			h.GetIncomingRequests(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	case "/api/friends/list":
		if r.Method == http.MethodGet {
			h.GetFriendsList(w, r)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}

	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
func (h *Handler) SendFriendRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		TargetID uint `json:"target_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.TargetID == 0 {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Проверка, не существует ли уже заявка
	var existing models.FriendRequest
	err := h.DB.Where("sender_id = ? AND request_id = ?", userID, req.TargetID).First(&existing).Error
	if err == nil {
		http.Error(w, "Request already exists", http.StatusConflict)
		return
	}

	request := models.FriendRequest{
		SenderID:  userID,
		RequestID: req.TargetID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	if err := h.DB.Create(&request).Error; err != nil {
		http.Error(w, "Failed to send request", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(request)
}

func (h *Handler) HandleFriendRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var req struct {
		RequestID uint   `json:"request_id"`
		Action    string `json:"action"` // "accept" or "decline"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var fr models.FriendRequest
	if err := h.DB.First(&fr, req.RequestID).Error; err != nil {
		http.Error(w, "Request not found", http.StatusNotFound)
		return
	}

	if fr.RequestID != userID {
		http.Error(w, "Not your request", http.StatusForbidden)
		return
	}

	if req.Action == "accept" {
		fr.Status = "accepted"
	} else {
		fr.Status = "declined"
	}

	if err := h.DB.Save(&fr).Error; err != nil {
		http.Error(w, "Failed to update", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fr)
}

func (h *Handler) GetIncomingRequests(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var requests []models.FriendRequest
	if err := h.DB.Preload("Sender").Where("request_id = ? AND status = 'pending'", userID).Find(&requests).Error; err != nil {
		http.Error(w, "Error fetching", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(requests)
}

func (h *Handler) GetFriendsList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	var friendIDs []uint
	h.DB.Raw(`
		SELECT CASE 
			WHEN sender_id = ? THEN request_id 
			WHEN request_id = ? THEN sender_id 
		END as friend_id
		FROM friends_requests
		WHERE (sender_id = ? OR request_id = ?) AND status = 'accepted'
	`, userID, userID, userID, userID).Scan(&friendIDs)

	var friends []models.User
	if len(friendIDs) > 0 {
		h.DB.Where("id IN ?", friendIDs).Find(&friends)
	}

	json.NewEncoder(w).Encode(friends)
}
