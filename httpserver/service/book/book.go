package book

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"github.com/storyofhis/books-management/httpserver/service"
	"gorm.io/gorm"
)

type bookSvc struct {
	repo repository.BookRepo
}

// CreateBook implements service.BookSvc.
func (svc *bookSvc) CreateBook(ctx context.Context, book *params.CreateBook, id uuid.UUID) *views.Response {
	param := models.Book{
		UserId:   id,
		AuthorId: book.AuthorId,
		Title:    book.Title,
		Isbn:     book.Isbn,
	}
	err := svc.repo.CreateBook(ctx, &param)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Book{
		Id:        param.Id,
		UserId:    param.UserId,
		AuthorId:  param.AuthorId,
		Title:     param.Title,
		Isbn:      param.Isbn,
		CreatedAt: param.CreatedAt,
		UpdatedAt: param.UpdatedAt,
	})
}

// DeleteBook implements service.BookSvc.
func (svc *bookSvc) DeleteBook(ctx context.Context, id uuid.UUID) *views.Response {
	book, err := svc.repo.GetBookById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	err = svc.repo.DeleteBook(ctx, id)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusNoContent, views.M_AUTHOR_SUCCESSFULLY_DELETED, views.Book{
		UserId: book.UserId,
	})
}

// GetAuthorById implements service.BookSvc.
func (svc *bookSvc) GetBookById(ctx context.Context, id uuid.UUID) *views.Response {
	book, err := svc.repo.GetBookById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_AUTHOR_NOT_FOUND, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, views.Book{
		Id:        book.Id,
		UserId:    book.UserId,
		AuthorId:  book.AuthorId,
		Title:     book.Title,
		Isbn:      book.Isbn,
		CreatedAt: book.CreatedAt,
		UpdatedAt: book.UpdatedAt,
	})
}

// GetBooks implements service.BookSvc.
func (svc *bookSvc) GetBooks(ctx context.Context) *views.Response {
	book, err := svc.repo.GetBooks(ctx)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	books := make([]views.Book, 0)
	for _, b := range book {
		books = append(books, views.Book{
			Id:        b.Id,
			UserId:    b.UserId,
			AuthorId:  b.AuthorId,
			Title:     b.Title,
			Isbn:      b.Isbn,
			CreatedAt: b.CreatedAt,
			UpdatedAt: b.UpdatedAt,
		})
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, books)
}

// UpdateAuthor implements service.BookSvc.
func (svc *bookSvc) UpdateBook(ctx context.Context, book *params.UpdateBook, id uuid.UUID) *views.Response {
	b, err := svc.repo.GetBookById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	b.AuthorId = book.AuthorId
	b.Title = book.Title
	b.Isbn = book.Isbn

	err = svc.repo.UpdateBook(ctx, b, id)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.UpdateBook{
		Id:        b.Id,
		UserId:    b.UserId,
		AuthorId:  b.AuthorId,
		Title:     b.Title,
		Isbn:      b.Isbn,
		UpdatedAt: b.UpdatedAt,
	})
}

func NewBookSvc(repo repository.BookRepo) service.BookSvc {
	return &bookSvc{
		repo: repo,
	}
}
