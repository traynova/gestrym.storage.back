package models

import "gorm.io/gorm"

const (
	UserTokenTypeActivation       = "activation"
	UserTokenTypePasswordRecovery = "password_recovery"
)

type UserTokenType struct {
	gorm.Model

	ID   uint   `gorm:"primaryKey" json:"id"`
	Type string `gorm:"type:varchar(50);uniqueIndex;not null" json:"type"`
}
