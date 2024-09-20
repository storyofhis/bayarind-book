package book_test

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
	"github.com/storyofhis/books-management/httpserver/service/book"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type bookSvcTest struct {
	repo    *repository.MockBookRepo
	service service.BookSvc
}

func newBookSvcTestTest(t *testing.T) bookSvcTest {
	mockRepo := repository.NewMockBookRepo(t)
	bookSvc := book.NewBookSvc(mockRepo)
	return bookSvcTest{
		repo:    mockRepo,
		service: bookSvc,
	}
}

func TestBookSvc_CreateBook(t *testing.T) {
	t.Run("success - it should return nil", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		instance.repo.EXPECT().CreateBook(mock.Anything, mock.Anything).Return(nil)
		res := instance.service.CreateBook(context.Background(), &params.CreateBook{}, uuid.New())
		assert.Equal(t, http.StatusCreated, res.Status)
	})

	t.Run("error - it should return an error if CreateBook returns an error", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		instance.repo.EXPECT().CreateBook(mock.Anything, mock.Anything).Return(assert.AnError)
		resp := instance.service.CreateBook(context.Background(), &params.CreateBook{}, uuid.New())
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
	})
}

func TestBookSvc_DeleteBook(t *testing.T) {
	t.Run("success - it should return no content", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()
		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(&models.Book{
			Id:       id,
			UserId:   uuid.New(),
			AuthorId: uuid.New(),
			Title:    "Test Book",
			Isbn:     "123456789",
		}, nil)

		instance.repo.EXPECT().DeleteBook(mock.Anything, id).Return(nil)
		res := instance.service.DeleteBook(context.Background(), id)
		assert.Equal(t, http.StatusNoContent, res.Status)
		assert.Nil(t, res.Payload)
	})

	t.Run("error - it should return an error if GetBookById returns an error", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)
		res := instance.service.DeleteBook(context.Background(), id)
		assert.Equal(t, http.StatusBadRequest, res.Status)
	})

	t.Run("error - it should return an error if DeleteBook fails", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()
		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(&models.Book{
			Id:       id,
			UserId:   uuid.New(),
			AuthorId: uuid.New(),
			Title:    "Test Book",
			Isbn:     "123456789",
		}, nil)

		instance.repo.EXPECT().DeleteBook(mock.Anything, id).Return(assert.AnError)
		res := instance.service.DeleteBook(context.Background(), id)
		assert.Equal(t, http.StatusInternalServerError, res.Status)
	})
}

func TestGetBookById(t *testing.T) {
	t.Run("success - it should return book details", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		mockBook := &models.Book{
			Id:        id,
			UserId:    uuid.New(),
			AuthorId:  uuid.New(),
			Title:     "Test Book",
			Isbn:      "123456789",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(mockBook, nil)
		res := instance.service.GetBookById(context.Background(), id)

		assert.Equal(t, http.StatusOK, res.Status)
		bookData, ok := res.Payload.(views.Book)
		assert.True(t, ok)
		assert.Equal(t, mockBook.Id, bookData.Id)
		assert.Equal(t, mockBook.Title, bookData.Title)
		assert.Equal(t, mockBook.Isbn, bookData.Isbn)
	})

	t.Run("error - it should return 400 if book not found", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)
		res := instance.service.GetBookById(context.Background(), id)
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_AUTHOR_NOT_FOUND, res.Message)
	})

	t.Run("error - it should return 500 if there is a database error", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(nil, assert.AnError)
		res := instance.service.GetBookById(context.Background(), id)
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}

func TestGetBooks(t *testing.T) {
	t.Run("success - it should return a list of books", func(t *testing.T) {
		instance := newBookSvcTestTest(t)

		// Create mock books as a slice of pointers
		mockBooks := []*models.Book{
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				AuthorId:  uuid.New(),
				Title:     "Test Book 1",
				Isbn:      "1234567890",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				Id:        uuid.New(),
				UserId:    uuid.New(),
				AuthorId:  uuid.New(),
				Title:     "Test Book 2",
				Isbn:      "0987654321",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		// Mock GetBooks to return the list of books
		instance.repo.EXPECT().GetBooks(mock.Anything).Return(mockBooks, nil)

		// Call GetBooks service
		res := instance.service.GetBooks(context.Background())

		// Assert response status is 200 OK
		assert.Equal(t, http.StatusOK, res.Status)

		// Assert response data contains the correct book details
		books, ok := res.Payload.([]views.Book)
		assert.True(t, ok)
		assert.Len(t, books, 2)

		// Assert the details of each book
		assert.Equal(t, mockBooks[0].Id, books[0].Id)
		assert.Equal(t, mockBooks[0].Title, books[0].Title)
		assert.Equal(t, mockBooks[1].Id, books[1].Id)
		assert.Equal(t, mockBooks[1].Title, books[1].Title)
	})

	t.Run("error - it should return 500 if there is a database error", func(t *testing.T) {
		instance := newBookSvcTestTest(t)

		// Mock GetBooks to return a generic error
		instance.repo.EXPECT().GetBooks(mock.Anything).Return(nil, assert.AnError)

		// Call GetBooks service
		res := instance.service.GetBooks(context.Background())

		// Assert response status is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}

func TestUpdateBook(t *testing.T) {
	t.Run("success - it should update the book details", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		// Create a mock book object before update
		mockBook := &models.Book{
			Id:        id,
			UserId:    uuid.New(),
			AuthorId:  uuid.New(),
			Title:     "Original Title",
			Isbn:      "1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Create the book update parameters
		updateParams := &params.UpdateBook{
			AuthorId: uuid.New(),
			Title:    "Updated Title",
			Isbn:     "0987654321",
		}

		// Mock GetBookById to return the original book
		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(mockBook, nil)

		// Mock UpdateBook to simulate successful update
		instance.repo.EXPECT().UpdateBook(mock.Anything, mockBook, id).Return(nil)

		// Call UpdateBook service
		res := instance.service.UpdateBook(context.Background(), updateParams, id)

		// Assert response status is 200 OK
		assert.Equal(t, http.StatusOK, res.Status)

		// Assert the book details have been updated in the response
		updatedBook, ok := res.Payload.(views.UpdateBook)
		assert.True(t, ok)
		assert.Equal(t, updateParams.Title, updatedBook.Title)
		assert.Equal(t, updateParams.Isbn, updatedBook.Isbn)
		assert.Equal(t, updateParams.AuthorId, updatedBook.AuthorId)
	})

	t.Run("error - it should return 404 if book not found", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		// Mock GetBookById to return a record not found error
		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(nil, gorm.ErrRecordNotFound)

		// Call UpdateBook service
		res := instance.service.UpdateBook(context.Background(), &params.UpdateBook{}, id)

		// Assert response status is 400 Bad Request
		assert.Equal(t, http.StatusBadRequest, res.Status)
		assert.Equal(t, views.M_BAD_REQUEST, res.Message)
	})

	t.Run("error - it should return 500 if there is a database error during update", func(t *testing.T) {
		instance := newBookSvcTestTest(t)
		id := uuid.New()

		// Create a mock book object before update
		mockBook := &models.Book{
			Id:        id,
			UserId:    uuid.New(),
			AuthorId:  uuid.New(),
			Title:     "Original Title",
			Isbn:      "1234567890",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock GetBookById to return the original book
		instance.repo.EXPECT().GetBookById(mock.Anything, id).Return(mockBook, nil)

		// Mock UpdateBook to return a database error
		instance.repo.EXPECT().UpdateBook(mock.Anything, mockBook, id).Return(assert.AnError)

		// Call UpdateBook service
		res := instance.service.UpdateBook(context.Background(), &params.UpdateBook{}, id)

		// Assert response status is 500 Internal Server Error
		assert.Equal(t, http.StatusInternalServerError, res.Status)
		assert.Equal(t, views.M_INTERNAL_SERVER_ERROR, res.Message)
	})
}
