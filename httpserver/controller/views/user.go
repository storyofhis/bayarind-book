package views

import (
	"time"

	"github.com/google/uuid"
)

type Register struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Token    string    `json:"token"`
}
