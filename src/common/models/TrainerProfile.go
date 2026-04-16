package models

import "gorm.io/gorm"

type TrainerProfile struct {
	gorm.Model

	PrimaryColor   *string `json:"primary_color"`
	SecondaryColor *string `json:"secondary_color"`
	ReferralCode   *string `gorm:"type:varchar(100)" json:"referral_code"`

	UserID  uint  `json:"user_id"`
	FilesID *uint `json:"files_id"`
	GimID   *uint `json:"gim_id"`

	User           User            `gorm:"foreignKey:UserID" json:"user"`
	Files          *Files          `gorm:"foreignKey:FilesID" json:"files"`
	GymProfile     *GymProfile     `gorm:"foreignKey:GimID;references:UserID" json:"gym_profile"`
	TrainerClients []TrainerClient `gorm:"foreignKey:TrainerProfileID" json:"trainer_clients"`
}
