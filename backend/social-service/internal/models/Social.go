package models

import "time"

type Challanges struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CreatorID       uint      `gorm:"not null" json:"creator_id"`
	Title           string    `gorm:"not null" json:"title"`
	Description     string    `gorm:"not null" json:"description"`
	ChallangeTypeID uint      `gorm:"not null" json:"challange_type_id"`
	GoalValue       int       `gorm:"not null" json:"goal_value"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`

	ChallangesType ChallangesTypes `gorm:"foreignKey:ChallangeTypeID"`
}

func (Challanges) TableName() string {
	return "challanges"
}

type ChallangesTypes struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
}

func (ChallangesTypes) TableName() string {
	return "challanges_types"
}

type ChallangeParticipants struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ChallangeTypeID uint      `gorm:"not null" json:"challange_type_id"`
	UserID          int       `gorm:"not null" json:"user_id"`
	Progress        int       `gorm:"not null" json:"progress"`
	Status          string    `gorm:"type:text;not null;check:status IN ('began','in_progress','achieved')" json:"status"`
	JoinAt          time.Time `gorm:"autoCreateTime" json:"join_at"`

	ChallangesType ChallangesTypes `gorm:"foreignKey:ChallangeTypeID"`
}

func (ChallangeParticipants) TableName() string {
	return "challange_participants"
}

type Messages struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SenderID   uint      `gorm:"not null" json:"sender_id"`
	RecieverID uint      `gorm:"not null" json:"reciever_id"`
	Content    string    `gorm:"not null" json:"content"`
	SentAt     time.Time `gorm:"autoCreateTime" json:"sent_at"`
}

func (Messages) TableName() string {
	return "messages"
}

type ActivityFeed struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	UserID    int                    `gorm:"not null" json:"user_id"`
	EventType string                 `gorm:"not null" json:"event_type"`
	EventData map[string]interface{} `gorm:"type:jsonb" json:"event_data"`
	CreatedAt time.Time              `gorm:"default:now()" json:"created_at"`
}

func (ActivityFeed) TableName() string {
	return "activity_feed"
}
