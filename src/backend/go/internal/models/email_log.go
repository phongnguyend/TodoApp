package models

import "time"

// EmailLog records every outbound email attempt.
// Maps to the email_logs table as defined in backend/README.md.
type EmailLog struct {
	ID           uint       `gorm:"primarykey;autoIncrement"              json:"id"`
	Recipient    string     `gorm:"type:varchar(255);not null"            json:"recipient"`
	Subject      string     `gorm:"type:varchar(500);not null"            json:"subject"`
	Body         string     `gorm:"type:text;not null"                    json:"body"`
	Status       string     `gorm:"type:varchar(50);not null;default:'pending'" json:"status"`
	CreatedAt    time.Time  `gorm:"autoCreateTime"                        json:"createdAt"`
	SentAt       *time.Time `gorm:"column:sent_at"                        json:"sentAt"`
	ErrorMessage *string    `gorm:"type:text;column:error_message"        json:"errorMessage"`
}

// TableName overrides the default table name.
func (EmailLog) TableName() string {
	return "email_logs"
}
