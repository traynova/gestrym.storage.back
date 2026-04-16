package models

import "gorm.io/gorm"

// Relación de muchos a muchos para Roles y Permisos (opcional/planeado)
type RolePermission struct {
	gorm.Model

	RoleID       uint       `gorm:"primaryKey"`
	PermissionID uint       `gorm:"primaryKey"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}
