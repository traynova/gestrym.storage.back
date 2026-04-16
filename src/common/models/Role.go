package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model

	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string `json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}
