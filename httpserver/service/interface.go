package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
)

type UserSvc interface {
	Register(ctx context.Context, user *params.Register) *views.Response
	Login(ctx context.Context, user *params.Login) *views.Response
}

type AuthorSvc interface {
	CreateAuthor(ctx context.Context, author *params.CreateAuthors, id uuid.UUID) *views.Response
	GetAuthors(ctx context.Context) *views.Response
	GetAuthorById(ctx context.Context, id uuid.UUID) *views.Response
	UpdateAuthor(ctx context.Context, author *params.UpdateAuthors, id uuid.UUID) *views.Response
	DeleteAuthor(ctx context.Context, id uuid.UUID) *views.Response
}

type BookSvc interface {
	CreateBook(ctx context.Context, book *params.CreateBook, id uuid.UUID) *views.Response
	GetBooks(ctx context.Context) *views.Response
	GetBookById(ctx context.Context, id uuid.UUID) *views.Response
	UpdateBook(ctx context.Context, book *params.UpdateBook, id uuid.UUID) *views.Response
	DeleteBook(ctx context.Context, id uuid.UUID) *views.Response
}
