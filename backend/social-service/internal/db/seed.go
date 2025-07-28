package db

import (
	"github.com/slickip/Healthy-summer-app/backend/social-service/internal/models" // замените на ваш путь к модели ChallangesTypes
	"gorm.io/gorm"
)

func SeedChallengeTypes(db *gorm.DB) error {
	challengeTypes := []models.ChallangesTypes{
		{Name: "Step Count", Description: "Reach a target number of steps over a defined period."},
		{Name: "Workout Frequency", Description: "Complete a set number of workouts within a time frame."},
		{Name: "Active Minutes", Description: "Accumulate active minutes through physical activity each day."},
		{Name: "Consistency Streak", Description: "Stay active without missing a day for a certain number of days."},
		{Name: "Water Intake", Description: "Drink at least a set amount of water each day."},
		{Name: "Healthy Meals", Description: "Log and consume a target number of healthy meals daily."},
		{Name: "No Sugar Week", Description: "Avoid added sugars for a week or more."},
		{Name: "Meal Logging Streak", Description: "Log every meal without skipping for consecutive days."},
		{Name: "Mindfulness", Description: "Practice meditation or mindfulness for a set duration daily."},
		{Name: "Gratitude Journal", Description: "Write down one thing you’re grateful for each day."},
		{Name: "Friend Boost", Description: "Send motivational messages to friends daily."},
		{Name: "Photo Share", Description: "Share daily photos of activities, food, or workouts."},
		{Name: "Total Wellness", Description: "Meet combined goals for steps, water, meals, and activity each day."},
		{Name: "Team Challenge", Description: "Join forces with friends to reach a group wellness goal."},
		{Name: "Beat Your Past", Description: "Surpass your previous week's performance by a percentage."},
	}

	for _, c := range challengeTypes {
		if err := db.Create(&c).Error; err != nil {
			return err
		}
	}

	return nil
}
