package models

import "gorm.io/gorm"

type Action struct {
	gorm.Model

	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
}
