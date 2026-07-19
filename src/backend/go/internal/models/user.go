package models

import "time"

// User is an account that can manage its own profile and password.
type User struct {
	ID           uint       `gorm:"primarykey;autoIncrement" json:"id"`
	Username     string     `gorm:"type:varchar(50);not null;uniqueIndex:ux_users_username,collate:nocase" json:"username"`
	Email        string     `gorm:"type:varchar(255);not null;uniqueIndex:ux_users_email,collate:nocase" json:"email"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`
	IsActive     bool       `gorm:"not null" json:"isActive"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}

func (User) TableName() string { return "users" }
