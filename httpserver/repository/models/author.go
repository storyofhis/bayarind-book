package models

import (
	"time"

	"github.com/google/uuid"
)

type Author struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID
	User      User `gorm:"foreignKey:UserId"`
	Name      string
	Birthdate time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}
