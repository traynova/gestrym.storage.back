package models

import "gorm.io/gorm"

type File struct {
	gorm.Model
	FileName    string `gorm:"type:varchar(255);not null" json:"file_name"`
	ContentType string `gorm:"type:varchar(100);not null" json:"content_type"`
	Size        int64  `gorm:"not null" json:"size"`
	URL         string `gorm:"type:text;not null" json:"url"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
	CollectionID string `gorm:"type:varchar(100);index" json:"collection_id"` // UUID to group files
	Service      string `gorm:"type:varchar(100)" json:"service"`           // e.g., "users", "exercises"
}
