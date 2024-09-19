package params

import "github.com/google/uuid"

type CreateBook struct {
	Title    string    `json:"title" validate:"required"`
	Isbn     string    `json:"isbn" validate:"required"`
	AuthorId uuid.UUID `json:"author_id" validate:"required"`
}

type UpdateBook struct {
	Title    string    `json:"title" validate:"required"`
	Isbn     string    `json:"isbn" validate:"required"`
	AuthorId uuid.UUID `json:"author_id" validate:"required"`
}
