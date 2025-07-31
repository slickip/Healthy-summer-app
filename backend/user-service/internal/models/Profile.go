package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	Achievements []Achievement `gorm:"foreignKey:UserID" json:"achievements"`
	// Исходящие запросы в друзья
	SentFriendRequests []FriendRequest `gorm:"foreignKey:SenderID" json:"sent_friend_requests"`
	// Входящие запросы в друзья
	ReceivedFriendRequests []FriendRequest `gorm:"foreignKey:RequestID" json:"received_friend_requests"`
}

func (User) TableName() string {
	return "users"
}

type FriendRequest struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SenderID  uint      `gorm:"not null" json:"sender_id"`
	RequestID uint      `gorm:"not null" json:"request_id"`
	Status    string    `gorm:"type:text;check:status IN ('pending','accepted','declined')" json:"status"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Sender   User `gorm:"foreignKey:SenderID" json:"sender"`
	Receiver User `gorm:"foreignKey:RequestID" json:"receiver"`
}

func (FriendRequest) TableName() string {
	return "friends_requests"
}

type Achievement struct {
	ID          uint      `gorm:"primaryKey"`
	UserID      uint      `gorm:"not null"`
	Title       string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text"`
	AchievedAt  time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:UserID"`
}

func (Achievement) TableName() string {
	return "achievements"
}

type Friend struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint `gorm:"not null"`
	FriendID  uint `gorm:"not null"`
	CreatedAt time.Time
}

func (Friend) TableName() string {
	return "friends"
}
