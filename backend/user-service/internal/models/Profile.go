package models

import (
	"time"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"unique;not null"`
	PasswordHash string `gorm:"not null"`
	DisplayName  string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	Achievements []Achievement `gorm:"foreignKey:UserID"`
	// Исходящие запросы в друзья
	SentFriendRequests []FriendRequest `gorm:"foreignKey:SenderID"`
	// Входящие запросы в друзья
	ReceivedFriendRequests []FriendRequest `gorm:"foreignKey:RequestID"`
}

func (User) TableName() string {
	return "users"
}

type FriendRequest struct {
	ID        uint      `gorm:"primaryKey"`
	SenderID  uint      `gorm:"not null"`
	RequestID uint      `gorm:"not null"`
	Status    string    `gorm:"type:text;check:status IN ('pending','accepted','declined')"`
	CreatedAt time.Time `gorm:"autoCreateTime"`

	Sender   User `gorm:"foreignKey:SenderID"`
	Receiver User `gorm:"foreignKey:RequestID"`
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
