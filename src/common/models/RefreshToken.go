package models

import (
	"time"

	"gorm.io/gorm"
)

type RefreshToken struct {
	// Campos comunes
	gorm.Model
	CreatedBy *uint
	UpdatedBy *uint

	// Campos específicos para la tabla
	Key        string    `gorm:"unique;not null;index"`
	ExpiryDate time.Time `gorm:"not null"`

	// Relaciones
	UserID uint `gorm:"not null"`
	User   User
}
