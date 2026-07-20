package models

import "time"

// User is an account that can manage its own profile and password.
type User struct {
	ID              uint       `gorm:"primarykey;autoIncrement" json:"id"`
	Username        string     `gorm:"type:varchar(50);not null;uniqueIndex:ux_users_username,collate:nocase" json:"username"`
	Email           string     `gorm:"type:varchar(255);not null;uniqueIndex:ux_users_email,collate:nocase" json:"email"`
	PasswordHash    string     `gorm:"type:varchar(255);not null" json:"-"`
	IsActive        bool       `gorm:"not null" json:"isActive"`
	CreatedAt       time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	CreatedByUserID *uint      `gorm:"index" json:"createdByUserId"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	UpdatedByUserID *uint      `gorm:"index" json:"updatedByUserId"`
	CreatedByUser   *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CreatedByUserID" json:"-"`
	UpdatedByUser   *User      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UpdatedByUserID" json:"-"`
}

func (User) TableName() string { return "users" }
