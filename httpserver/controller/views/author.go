package views

import (
	"time"

	"github.com/google/uuid"
)

type UpdateAuthor struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateAuthor struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
	CreatedAt time.Time `json:"created_at"`
}

type Author struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Birthdate time.Time `json:"birthdate"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
