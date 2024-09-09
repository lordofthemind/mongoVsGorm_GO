package types

import (
	"time"

	"github.com/google/uuid"
)

// Author represents the author entity used with both MongoDB and GORM.
type Author struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey" bson:"_id"`                                          // UUID type for primary key
	Name        string     `json:"name" gorm:"column:name" bson:"name"`                                                // Name of the author
	Bio         string     `json:"bio,omitempty" gorm:"column:bio" bson:"bio,omitempty"`                               // Biography of the author, omitted if empty
	Email       string     `json:"email" gorm:"column:email" bson:"email"`                                             // Email address of the author
	DateOfBirth *time.Time `json:"date_of_birth,omitempty" gorm:"column:date_of_birth" bson:"date_of_birth,omitempty"` // Date of birth, omitted if nil
}
