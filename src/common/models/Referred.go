package models

type Referred struct {
	UserID       uint   `gorm:"not null" json:"user_id"`
	User         User   `gorm:"foreignKey:UserID" json:"user"`
	ReferralCode string `gorm:"type:varchar(100);uniqueIndex;not null" json:"referral_code"`
}
