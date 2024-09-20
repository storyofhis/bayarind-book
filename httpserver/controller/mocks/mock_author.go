package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/mock"
)

type MockAuthorSvc struct {
	mock.Mock
}

// DeleteAuthor implements service.AuthorSvc.
func (m *MockAuthorSvc) DeleteAuthor(ctx context.Context, id uuid.UUID) *views.Response {
	args := m.Called(ctx, id)
	return args.Get(0).(*views.Response)
}

// GetAuthorById implements service.AuthorSvc.
func (m *MockAuthorSvc) GetAuthorById(ctx context.Context, id uuid.UUID) *views.Response {
	args := m.Called(ctx, id)
	return args.Get(0).(*views.Response)
}

// GetAuthors implements service.AuthorSvc.
func (m *MockAuthorSvc) GetAuthors(ctx context.Context) *views.Response {
	args := m.Called(ctx)
	return args.Get(0).(*views.Response)
}

// UpdateAuthor implements service.AuthorSvc.
func (m *MockAuthorSvc) UpdateAuthor(ctx context.Context, author *params.UpdateAuthors, id uuid.UUID) *views.Response {
	args := m.Called(ctx, author, id)
	return args.Get(0).(*views.Response)
}

func (m *MockAuthorSvc) CreateAuthor(ctx context.Context, authorParams *params.CreateAuthors, userId uuid.UUID) *views.Response {
	args := m.Called(ctx, authorParams, userId)
	return args.Get(0).(*views.Response)
}
