package models

import (
	"time"

	"gorm.io/gorm"
)

type UserToken struct {
	gorm.Model

	ID              uint          `gorm:"primaryKey" json:"id"`
	UserID          uint          `gorm:"not null;index" json:"user_id"`
	User            User          `gorm:"foreignKey:UserID" json:"-"`
	UserTokenTypeID uint          `gorm:"not null;index" json:"user_token_type_id"`
	UserTokenType   UserTokenType `gorm:"foreignKey:UserTokenTypeID" json:"-"`
	Token           string        `gorm:"type:text;not null;uniqueIndex" json:"token"`
	IsRevoked       bool          `gorm:"default:false" json:"is_revoked"`
	ExpiresAt       time.Time     `gorm:"not null" json:"expires_at"`
}
