package views

import (
	"time"

	"github.com/google/uuid"
)

type CreateBook struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	AuthorId  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Isbn      string    `json:"isbn"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateBook struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	AuthorId  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Isbn      string    `json:"isbn"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Book struct {
	Id        uuid.UUID `json:"id"`
	UserId    uuid.UUID `json:"user_id"`
	AuthorId  uuid.UUID `json:"author_id"`
	Title     string    `json:"title"`
	Isbn      string    `json:"isbn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
