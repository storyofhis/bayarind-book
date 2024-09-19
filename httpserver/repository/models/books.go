package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID
	User      User `gorm:"foreignKey:UserId"`
	AuthorId  uuid.UUID
	Author    Author `gorm:"foreignKey:AuthorId"`
	Title     string
	Isbn      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
