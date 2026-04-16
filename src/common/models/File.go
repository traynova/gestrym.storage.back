package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName    string `gorm:"type:varchar(255);not null" json:"file_name"`
	ContentType string `gorm:"type:varchar(100);not null" json:"content_type"`
	Size        int64  `gorm:"not null" json:"size"`
	URL         string `gorm:"type:text;not null" json:"url"`
	Collection  string `gorm:"type:varchar(100)" json:"collection"`   // e.g., "exercise-images"
	EntityID    string `gorm:"type:varchar(100)" json:"entity_id"`    // e.g., exercise ID
	EntityType  string `gorm:"type:varchar(100)" json:"entity_type"`  // e.g., "exercise"
}
