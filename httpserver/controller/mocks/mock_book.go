// booksvc_mock.go (inside your test folder)
package mocks

import (
	"context"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/mock"
)

type MockBookSvc struct {
	mock.Mock
}

func (m *MockBookSvc) CreateBook(ctx context.Context, bookParams *params.CreateBook, userId uuid.UUID) *views.Response {
	args := m.Called(ctx, bookParams, userId)
	return args.Get(0).(*views.Response)
}

func (m *MockBookSvc) GetBooks(ctx context.Context) *views.Response {
	args := m.Called(ctx)
	return args.Get(0).(*views.Response)
}

func (m *MockBookSvc) GetBookById(ctx context.Context, id uuid.UUID) *views.Response {
	args := m.Called(ctx, id)
	return args.Get(0).(*views.Response)
}

func (m *MockBookSvc) UpdateBook(ctx context.Context, bookParams *params.UpdateBook, id uuid.UUID) *views.Response {
	args := m.Called(ctx, bookParams, id)
	return args.Get(0).(*views.Response)
}

func (m *MockBookSvc) DeleteBook(ctx context.Context, id uuid.UUID) *views.Response {
	args := m.Called(ctx, id)
	return args.Get(0).(*views.Response)
}
