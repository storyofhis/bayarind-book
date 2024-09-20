package author_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"github.com/storyofhis/books-management/httpserver/service"
	"github.com/storyofhis/books-management/httpserver/service/author"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type authorSvcTest struct {
	repo    *repository.MockAuthorRepo
	service service.AuthorSvc
}

func newAuthorSvcTest(t *testing.T) authorSvcTest {
	mockRepo := repository.NewMockAuthorRepo(t)
	authorSvc := author.NewAuthorSvc(mockRepo)
	return authorSvcTest{
		repo:    mockRepo,
		service: authorSvc,
	}
}

func TestAuthorSvc_CreateAuthor(t *testing.T) {
	t.Run("success - it shoult return nil", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		instance.repo.EXPECT().CreateAuthor(mock.Anything, mock.Anything).Return(nil)
		res := instance.service.CreateAuthor(context.Background(), &params.CreateAuthors{}, uuid.New())
		assert.Equal(t, http.StatusCreated, res.Status)
	})

	t.Run("error - it should return an error if CreateAuthor returns an error", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		instance.repo.EXPECT().CreateAuthor(mock.Anything, mock.Anything).Return(assert.AnError)
		resp := instance.service.CreateAuthor(context.Background(), &params.CreateAuthors{}, uuid.New())
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
	})
}

func TestAuthorSvc_DeleteAuthor(t *testing.T) {
	t.Run("success - it should delete the author", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()
		mockAuthor := &models.Author{
			Id:     id,
			UserId: uuid.New(),
		}

		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(mockAuthor, nil)
		instance.repo.EXPECT().DeleteAuthor(mock.Anything, id).Return(nil)
		res := instance.service.DeleteAuthor(context.Background(), id)

		assert.Equal(t, http.StatusNoContent, res.Status)
		authorData, ok := res.Payload.(views.Author)
		assert.True(t, ok)
		assert.Equal(t, mockAuthor.UserId, authorData.UserId)
	})

	t.Run("error - it should return 404 if author not found", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()
		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

		res := instance.service.DeleteAuthor(context.Background(), id)
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_BAD_REQUEST, res.Message)
	})

	t.Run("error - it should return 500 if there is a database error during delete", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()

		mockAuthor := &models.Author{
			Id:     id,
			UserId: uuid.New(),
		}
		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(mockAuthor, nil)
		instance.repo.EXPECT().DeleteAuthor(mock.Anything, id).Return(assert.AnError)
		res := instance.service.DeleteAuthor(context.Background(), id)

		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}

func TestAuthorSvc_GetAuthorById(t *testing.T) {
	t.Run("success - it should return author details", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()

		mockAuthor := &models.Author{
			Id:        id,
			UserId:    uuid.New(),
			Name:      "John Doe",
			Birthdate: time.Date(1980, time.January, 1, 0, 0, 0, 0, time.UTC),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(mockAuthor, nil)

		res := instance.service.GetAuthorById(context.Background(), id)

		assert.Equal(t, http.StatusOK, res.Status)
		authorData, ok := res.Payload.(views.Author)
		assert.True(t, ok)
		assert.Equal(t, mockAuthor.Id, authorData.Id)
		assert.Equal(t, mockAuthor.Name, authorData.Name)
		assert.Equal(t, mockAuthor.Birthdate, authorData.Birthdate)
	})

	t.Run("error - it should return 404 if author not found", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()
		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

		res := instance.service.GetAuthorById(context.Background(), id)

		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_AUTHOR_NOT_FOUND, res.Message)
	})

	t.Run("error - it should return 500 if there is a database error", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()

		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(nil, assert.AnError)
		res := instance.service.GetAuthorById(context.Background(), id)
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}

func TestAuthorSvc_GetAuthors(t *testing.T) {
	t.Run("success - it should return a list of authors", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		mockAuthors := []*models.Author{
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				Name:      "John Doe",
				Birthdate: time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				Name:      "Jane Smith",
				Birthdate: time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		instance.repo.EXPECT().GetAuthors(mock.Anything).Return(mockAuthors, nil)
		res := instance.service.GetAuthors(context.Background())
		assert.Equal(t, http.StatusOK, res.Status)

		authorsData, ok := res.Payload.([]views.Author)
		assert.True(t, ok)
		assert.Len(t, authorsData, len(mockAuthors))
	})

	t.Run("error - it should return 500 if there is a database error", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		instance.repo.EXPECT().GetAuthors(mock.Anything).Return(nil, assert.AnError)
		res := instance.service.GetAuthors(context.Background())

		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}

func TestAuthorSvc_UpdateAuthor(t *testing.T) {
	t.Run("success - it should update the author and return the updated details", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()

		mockAuthor := &models.Author{
			Id:        id,
			UserId:    uuid.New(),
			Name:      "John Doe",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updatedAuthor := &params.UpdateAuthors{
			Name:      "John Updated",
			Birthdate: time.Now(),
		}

		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(mockAuthor, nil)
		instance.repo.EXPECT().UpdateAuthor(mock.Anything, mock.Anything, id).Return(nil)
		res := instance.service.UpdateAuthor(context.Background(), updatedAuthor, id)

		assert.Equal(t, http.StatusOK, res.Status)
		updatedAuthorData, ok := res.Payload.(views.UpdateAuthor)
		assert.True(t, ok)
		assert.Equal(t, updatedAuthor.Name, updatedAuthorData.Name)
		assert.Equal(t, updatedAuthor.Birthdate, updatedAuthorData.Birthdate)
	})

	t.Run("error - it should return 400 if author not found", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()
		updatedAuthor := &params.UpdateAuthors{
			Name:      "John Updated",
			Birthdate: time.Now(),
		}

		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)
		res := instance.service.UpdateAuthor(context.Background(), updatedAuthor, id)

		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_BAD_REQUEST, res.Message)
	})

	t.Run("error - it should return 500 if there is an error during update", func(t *testing.T) {
		instance := newAuthorSvcTest(t)
		id := uuid.New()

		mockAuthor := &models.Author{
			Id:        id,
			UserId:    uuid.New(),
			Name:      "John Doe",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		updatedAuthor := &params.UpdateAuthors{
			Name:      "John Updated",
			Birthdate: time.Now(),
		}

		instance.repo.EXPECT().GetAuthorById(mock.Anything, id).Return(mockAuthor, nil)
		instance.repo.EXPECT().UpdateAuthor(mock.Anything, mock.Anything, id).Return(assert.AnError)

		res := instance.service.UpdateAuthor(context.Background(), updatedAuthor, id)
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}
