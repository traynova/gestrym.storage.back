package models

import "gorm.io/gorm"

type User struct {
	gorm.Model

	Email             string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password          string `gorm:"type:varchar(255)" json:"password"`
	FullName          string `gorm:"type:varchar(100);not null" json:"name"`
	Prefix            string `json:"prefix" binding:"required"`
	Phone             string `gorm:"type:varchar(20)" json:"phone"`
	IsActive          bool   `gorm:"default:true" json:"is_active"`
	AccessFailedCount int    `gorm:"default:0" json:"access_failed_count"`
	LoguinMethodId    *uint  `json:"login_method_id"`
	EmailConfirmed    bool   `gorm:"default:false" json:"email_confirmed"`

	RoleID uint `gorm:"not null" json:"role_id"`

	Role Role `gorm:"foreignKey:RoleID" json:"role"`
}
