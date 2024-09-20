package mocks

import (
	"context"

	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/mock"
)

type MockUserSvc struct {
	mock.Mock
}

// Register mocks the Register function of the UserSvc
func (m *MockUserSvc) Register(ctx context.Context, user *params.Register) *views.Response {
	args := m.Called(ctx, user)
	return args.Get(0).(*views.Response)
}

// Login mocks the Login function of the UserSvc
func (m *MockUserSvc) Login(ctx context.Context, user *params.Login) *views.Response {
	args := m.Called(ctx, user)
	return args.Get(0).(*views.Response)
}
