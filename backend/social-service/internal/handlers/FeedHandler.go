package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/middleware"
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/models"
	"gorm.io/gorm"
)

type FeedHandler struct {
	DB *gorm.DB
}

// FriendsFeed godoc
// @Summary Get friends activity feed
// @Description Returns recent activities of the authenticated user's friends
// @Tags social
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {array} models.ActivityFeed "List of friends' activities, empty array if no friends or activities"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /api/social/feed/friends [get]
func (h *ChallengeHandler) FriendsFeed(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userIDValue := ctx.Value(middleware.ContextUserIDKey)
	if userIDValue == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := userIDValue.(uint)

	// Получить список ID друзей через friend_requests
	var friendIDs []uint
	h.DB.Raw(`
		SELECT CASE 
			WHEN sender_id = ? THEN receiver_id 
			WHEN receiver_id = ? THEN sender_id 
		END as friend_id
		FROM friend_requests
		WHERE (sender_id = ? OR receiver_id = ?) AND status = 'accepted'`, userID, userID, userID, userID).Scan(&friendIDs)

	if len(friendIDs) == 0 {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode([]any{})
		return
	}

	// Получаем публичную активность друзей
	var feed []models.ActivityFeed
	if err := h.DB.Where("user_id IN ?", friendIDs).Order("created_at DESC").Limit(100).Find(&feed).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
