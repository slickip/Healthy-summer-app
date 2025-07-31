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

	case "/api/users/search":
		if r.Method == http.MethodGet {
			h.SearchAllUsers(w, r)
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
	type FriendRequestResponse struct {
		ID        uint      `json:"id"`
		SenderID  uint      `json:"sender_id"`
		RequestID uint      `json:"request_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	resp := FriendRequestResponse{
		ID:        request.ID,
		SenderID:  request.SenderID,
		RequestID: request.RequestID,
		Status:    request.Status,
		CreatedAt: request.CreatedAt,
	}
	json.NewEncoder(w).Encode(resp)
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
	type FriendRequestResponse struct {
		ID        uint      `json:"id"`
		SenderID  uint      `json:"sender_id"`
		RequestID uint      `json:"request_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	resp := FriendRequestResponse{
		ID:        fr.ID,
		SenderID:  fr.SenderID,
		RequestID: fr.RequestID,
		Status:    fr.Status,
		CreatedAt: fr.CreatedAt,
	}
	json.NewEncoder(w).Encode(resp)
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
	type FriendRequestResponse struct {
		ID        uint      `json:"id"`
		SenderID  uint      `json:"sender_id"`
		RequestID uint      `json:"request_id"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"created_at"`
	}
	var resp []FriendRequestResponse
	for _, fr := range requests {
		resp = append(resp, FriendRequestResponse{
			ID:        fr.ID,
			SenderID:  fr.SenderID,
			RequestID: fr.RequestID,
			Status:    fr.Status,
			CreatedAt: fr.CreatedAt,
		})
	}
	json.NewEncoder(w).Encode(resp)
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
	type UserResponse struct {
		ID          uint      `json:"id"`
		Email       string    `json:"email"`
		DisplayName string    `json:"display_name"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}
	var resp []UserResponse
	for _, u := range friends {
		resp = append(resp, UserResponse{
			ID:          u.ID,
			Email:       u.Email,
			DisplayName: u.DisplayName,
			CreatedAt:   u.CreatedAt,
			UpdatedAt:   u.UpdatedAt,
		})
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if query == "" || userIDValue == nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	userID := userIDValue.(uint)

	var users []models.User
	subQuery := h.DB.
		Table("friends_requests").
		Select("request_id").
		Where("sender_id = ? AND status IN ('pending', 'accepted')", userID)

	h.DB.
		Where("id != ? AND id NOT IN (?)", userID, subQuery).
		Where("display_name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(10).Find(&users)

	type UserSearchResponse struct {
		ID          uint   `json:"id"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	var resp []UserSearchResponse
	for _, u := range users {
		resp = append(resp, UserSearchResponse{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
		})
	}

	json.NewEncoder(w).Encode(resp)
}

func (h *Handler) SearchAllUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var users []models.User
	h.DB.
		Where("display_name ILIKE ? OR email ILIKE ?", "%"+query+"%", "%"+query+"%").
		Limit(20).Find(&users)

	type UserSearchResponse struct {
		ID          uint   `json:"id"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	var resp []UserSearchResponse
	for _, u := range users {
		resp = append(resp, UserSearchResponse{
			ID:          u.ID,
			DisplayName: u.DisplayName,
			Email:       u.Email,
		})
	}

	json.NewEncoder(w).Encode(resp)
}
