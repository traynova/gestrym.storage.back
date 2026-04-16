package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model

	ID          uint        `gorm:"primaryKey" json:"id"`
	RoleID      uint        `gorm:"not null" json:"role_id"`
	ActionID    uint        `gorm:"not null" json:"action_id"`
	ResourceID  uint        `gorm:"not null" json:"resource_id"`
	Role        Role        `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Action      Action      `gorm:"foreignKey:ActionID" json:"action,omitempty"`
	AccessLevel AccessLevel `gorm:"foreignKey:ResourceID" json:"access_level,omitempty"`
	IsActive    bool        `gorm:"default:true" json:"is_active"`
}
