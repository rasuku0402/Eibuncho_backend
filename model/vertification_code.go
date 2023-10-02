package model

import (
	"time"
)

// VerificationCode represents the model for a verification code
type VerificationCode struct {
	PinID     uint       `json:"pinid" gorm:"primaryKey"` // Primary key
	UserName  string     `json:"username"`                // User name
	Email     string     `json:"email"`                   // User email
	Code      string     `json:"code"`                    // Generated PIN code
	Password  string     `json:"password"`                // Password hash
	ExpiresAt time.Time  `json:"expires_at"`              // Expiration time of the PIN code
	CreatedAt time.Time  `json:"created_at"`              // Timestamp of when the PIN code was generated
	UpdatedAt time.Time  `json:"updated_at"`              // Timestamp of when the PIN code was last updated
	DeletedAt *time.Time `gorm:"index"`                   // Soft deletion timestamp
}

// TableName sets the insert table name for this struct type
func (v *VerificationCode) TableName() string {
	return "verification_codes"
}

type VerificationRequest struct {
	Code string `json:"code"`
}

type Verificationresponse struct {
	UserName string `json:"username"` // User name
	Email    string `json:"email"`    // User email
	Password string `json:"password"` // Password hash
}
