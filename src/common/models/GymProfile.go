package models

import "gorm.io/gorm"

type GymProfile struct {
	gorm.Model

	City           string  `gorm:"type:varchar(100);not null" json:"city"`
	Department     string  `gorm:"type:varchar(100);not null" json:"department"`
	Country        string  `gorm:"type:varchar(100);not null" json:"country"`
	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	Workstation    *string `json:"workstation"`
	ReferralCode   *string `gorm:"type:varchar(100)" json:"referral_code"`

	UserID  uint  `gorm:"uniqueIndex" json:"user_id"`
	FilesID *uint `json:"files_id"`

	User  User   `gorm:"foreignKey:UserID" json:"user"`
	Files *Files `gorm:"foreignKey:FilesID" json:"files"`
}
