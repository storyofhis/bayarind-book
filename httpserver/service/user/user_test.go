package user_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"github.com/storyofhis/books-management/httpserver/service"
	"github.com/storyofhis/books-management/httpserver/service/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userSvcTest struct {
	repo    *repository.MockUserRepo
	service service.UserSvc
}

func newUserSvcTestTest(t *testing.T) userSvcTest {
	mockRepo := repository.NewMockUserRepo(t)
	userSvc := user.NewUserSvc(mockRepo)
	return userSvcTest{
		repo:    mockRepo,
		service: userSvc,
	}
}

func TestUserSvc_Register(t *testing.T) {
	t.Run("success - it should return nil", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
		instance.repo.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(nil)
		expected := views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Register{
			Username: "username",
		})
		res := instance.service.Register(context.Background(), &params.Register{
			Username: "username",
		})
		assert.Equal(t, http.StatusCreated, res.Status)
		assert.Equal(t, expected.Payload.(views.Register).Username, res.Payload.(views.Register).Username)
	})
	t.Run("error - it should return an error if CreateUser returns an error", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)
		instance.repo.EXPECT().CreateUser(mock.Anything, mock.Anything).Return(assert.AnError)
		resp := instance.service.Register(context.Background(), &params.Register{})
		assert.Equal(t, http.StatusInternalServerError, resp.Error)
	})
	t.Run("error - it should return an error if GetUserByUsername returns an error", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(nil, assert.AnError)
		resp := instance.service.Register(context.Background(), &params.Register{})
		assert.Equal(t, http.StatusInternalServerError, resp.Message)
	})
	t.Run("error - it should return M_USERNAME_ALREADY_USED if GetUserByUsername returns a user", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(&models.User{}, nil)
		resp := instance.service.Register(context.Background(), &params.Register{})
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Equal(t, views.M_USERNAME_ALREADY_USED, resp.Message)
	})
}

func TestUserSvc_Login(t *testing.T) {
	t.Run("success - it should return a valid token for correct credentials", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		// Mocking a valid user in the database
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		expectedUser := &models.User{
			Id:       uuid.New(),
			Username: "username",
			Password: string(hashedPassword),
		}

		// Mocking the repository responses
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(expectedUser, nil)

		// Calling the login service
		res := instance.service.Login(context.Background(), &params.Login{
			Username: "username",
			Password: "password",
		})

		// Asserting the results
		assert.Equal(t, http.StatusOK, res.Status)
		assert.NotEmpty(t, res.Payload.(views.Login).Token)
		assert.Equal(t, expectedUser.Username, res.Payload.(views.Login).Username)
	})

	t.Run("error - it should return an error if the username is not found", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		// Mock the repository to return ErrRecordNotFound for GetUserByUsername
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound)

		// Call the login service with non-existing username
		res := instance.service.Login(context.Background(), &params.Login{
			Username: "nonexistent",
			Password: "password",
		})

		// Assert that an invalid credentials error is returned
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_INVALID_CREDENTIALS, res.Message)
	})

	t.Run("error - it should return an error if the password is incorrect", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		// Mocking a valid user but with a different password
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
		expectedUser := &models.User{
			Id:       uuid.New(),
			Username: "username",
			Password: string(hashedPassword),
		}

		// Mocking the repository response
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(expectedUser, nil)

		// Call the login service with an incorrect password
		res := instance.service.Login(context.Background(), &params.Login{
			Username: "username",
			Password: "wrongpassword",
		})

		// Assert that an invalid credentials error is returned
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_INVALID_CREDENTIALS, res.Message)
	})

	t.Run("error - it should return an error if GetUserByUsername returns an error", func(t *testing.T) {
		instance := newUserSvcTestTest(t)
		// Mock the repository to return an arbitrary error
		instance.repo.EXPECT().GetUserByUsername(mock.Anything, mock.Anything).Return(nil, assert.AnError)

		// Call the login service with valid inputs
		res := instance.service.Login(context.Background(), &params.Login{
			Username: "username",
			Password: "password",
		})

		// Assert that an internal server error is returned
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}
