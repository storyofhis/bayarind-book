package params

import "time"

type CreateAuthors struct {
	Name      string    `json:"name" validate:"required"`
	Birthdate time.Time `json:"birthdate" validate:"required"`
}

type UpdateAuthors struct {
	Name      string    `json:"name" validate:"required"`
	Birthdate time.Time `json:"birthdate" validate:"required"`
	UpdateAt  time.Time `json:"updated_at"`
}
