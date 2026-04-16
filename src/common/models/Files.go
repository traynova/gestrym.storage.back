package models

import "gorm.io/gorm"

type Files struct {
	gorm.Model

	FileName  string `gorm:"type:varchar(255)" json:"file_name"`
	FilePath  string `gorm:"type:varchar(255)" json:"file_path"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
	FileSize  int64  `json:"file_size"`
	Extension string `gorm:"type:varchar(50)" json:"extension"`
}
