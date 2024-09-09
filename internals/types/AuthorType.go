package types

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey" bson:"_id"`
	Name        string     `json:"name" gorm:"column:name" bson:"name"`
	Bio         string     `json:"bio,omitempty" gorm:"column:bio" bson:"bio,omitempty"`
	Email       string     `json:"email" gorm:"column:email" bson:"email"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" bson:"date_of_birth,omitempty"`
}
