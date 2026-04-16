package models

import "gorm.io/gorm"

type GymClient struct {
	gorm.Model

	GymUserID uint `json:"gym_user_id"`
	ClientID  uint `json:"client_id"`

	Gym    *User `gorm:"foreignKey:GymUserID" json:"gym"`
	Client *User `gorm:"foreignKey:ClientID" json:"client"`
}
