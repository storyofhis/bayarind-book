package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/repository/models"
)

type UserRepo interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUsers(ctx context.Context, criteria map[string]interface{}) ([]*models.User, error)
	UpdateUserById(ctx context.Context, id uuid.UUID) error
}

type BookRepo interface {
	CreateBook(ctx context.Context, book *models.Book) error
	GetBooks(ctx context.Context) ([]*models.Book, error)
	GetBookById(ctx context.Context, id uuid.UUID) (*models.Book, error)
	DeleteBook(ctx context.Context, id uuid.UUID) error
	UpdateBook(ctx context.Context, book *models.Book, id uuid.UUID) error
}

type AuthorRepo interface {
	CreateAuthor(ctx context.Context, author *models.Author) error
	GetAuthors(ctx context.Context) ([]*models.Author, error)
	GetAuthorById(ctx context.Context, id uuid.UUID) (*models.Author, error)
	UpdateAuthor(ctx context.Context, author *models.Author, id uuid.UUID) error
	DeleteAuthor(ctx context.Context, id uuid.UUID) error
}
