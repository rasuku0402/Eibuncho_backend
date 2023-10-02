package model

import "time"

type Sentences struct {
	SentenceID   uint      `json:"sentenceid" gorm:"primaryKey"`
	User         User      `json:"user" gorm:"foreignKey:UserId; constraint:OnDelete:CASCADE"`
	UserId       uint      `json:"user_id" gorm:"not null"`
	Japanese_sen string    `json:"japanese_sen"`
	English_sen  string    `json:"english_sen"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SentenceResponse struct {
	SentenceID   uint      `json:"sentenceid" gorm:"primaryKey"`
	Japanese_sen string    `json:"japanese_sen"`
	English_sen  string    `json:"english_sen"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
