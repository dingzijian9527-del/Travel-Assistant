package repository

import (
	"time"

	aiagent "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/ai_agent"
	trip "github.com/dingzijian9527-del/Travel-Assistant/kitex_gen/trip"
)

type User struct {
	ID           string
	Phone        string
	PasswordHash string
	Nickname     string
	AvatarURL    string
	HomeCity     string
	CurrentCity  string
	Status       int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

type UserSettings struct {
	TripReminderEnabled          bool
	PriceReminderEnabled         bool
	PersonalizedRecommendEnabled bool
}

type Trip struct {
	ID             string
	UserID         int64
	Title          string
	Subtitle       string
	Destination    string
	DateRange      string
	DayCount       int32
	People         string
	BudgetLevel    string
	SourceQuestion string
	SourceReply    string
	Summary        *trip.TripSummary
	Days           []*trip.TripDay
	Budget         []*trip.TripBudget
	Alerts         []string
	Saved          bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      *time.Time
}

type ChatMessage struct {
	Role    aiagent.ChatRole
	Content string
}
