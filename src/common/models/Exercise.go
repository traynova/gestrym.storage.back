package models

import "time"

// Exercise represents a fitness exercise stored internally in the system.
type Exercise struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ExtID     string    `gorm:"size:100;uniqueIndex" json:"extId"` // The ID from external sources like ExerciseDB for deduplication
	Name      string    `gorm:"size:255;not null" json:"name"`
	BodyPart  string    `gorm:"size:100" json:"bodyPart"`
	Target    string    `gorm:"size:100" json:"target"`
	Equipment string    `gorm:"size:100" json:"equipment"`
	GifURL    string    `gorm:"size:500" json:"gifUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
