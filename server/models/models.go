package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex;size:50;not null"`
	Password  string    `json:"-" gorm:"size:255;not null"`
	Email     string    `json:"email" gorm:"size:100"`
	Region    string    `json:"region" gorm:"size:50;default:全国"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type EmissionFactor struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Category    string    `json:"category" gorm:"size:20;not null"`
	Item        string    `json:"item" gorm:"size:50;not null"`
	Factor      float64   `json:"factor" gorm:"type:decimal(10,4);not null"`
	Unit        string    `json:"unit" gorm:"size:20;not null"`
	Description string    `json:"description" gorm:"size:200"`
	CreatedAt   time.Time `json:"created_at"`
}

type CarbonRecord struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;index"`
	RecordDate string    `json:"record_date" gorm:"type:date;not null"`
	Category   string    `json:"category" gorm:"size:20;not null"`
	Item       string    `json:"item" gorm:"size:50;not null"`
	Amount     float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Emission   float64   `json:"emission" gorm:"type:decimal(10,4);not null"`
	Note       string    `json:"note" gorm:"size:200"`
	CreatedAt  time.Time `json:"created_at"`
}

type ReductionGoal struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"not null;index"`
	TargetEmission float64   `json:"target_emission" gorm:"type:decimal(10,2);not null"`
	StartDate      string    `json:"start_date" gorm:"type:date;not null"`
	EndDate        string    `json:"end_date" gorm:"type:date;not null"`
	Status         string    `json:"status" gorm:"size:20;default:active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegionalAverage struct {
	ID          uint    `json:"id" gorm:"primaryKey"`
	Region      string  `json:"region" gorm:"size:50;not null"`
	Month       string  `json:"month" gorm:"size:7;not null"`
	AvgEmission float64 `json:"avg_emission" gorm:"type:decimal(10,2);not null"`
}
