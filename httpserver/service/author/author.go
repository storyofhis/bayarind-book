package author

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

type authorSvc struct {
	repo repository.AuthorRepo
}

// CreateAuthor implements service.AuthorSvc.
func (svc *authorSvc) CreateAuthor(ctx context.Context, author *params.CreateAuthors, id uuid.UUID) *views.Response {
	param := models.Author{
		UserId:    id,
		Name:      author.Name,
		Birthdate: author.Birthdate,
	}

	err := svc.repo.CreateAuthor(ctx, &param)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Author{
		Id:        param.Id,
		UserId:    param.UserId,
		Name:      param.Name,
		Birthdate: param.Birthdate,
		CreatedAt: param.CreatedAt,
		UpdatedAt: param.UpdatedAt,
	})
}

// DeleteAuthor implements service.AuthorSvc.
func (svc *authorSvc) DeleteAuthor(ctx context.Context, id uuid.UUID) *views.Response {
	author, err := svc.repo.GetAuthorById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	err = svc.repo.DeleteAuthor(ctx, id)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusNoContent, views.M_AUTHOR_SUCCESSFULLY_DELETED, views.Author{
		UserId: author.UserId,
	})
}

// GetAuthorById implements service.AuthorSvc.
func (svc *authorSvc) GetAuthorById(ctx context.Context, id uuid.UUID) *views.Response {
	author, err := svc.repo.GetAuthorById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_AUTHOR_NOT_FOUND, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.Author{
		Id:        author.Id,
		UserId:    author.UserId,
		Name:      author.Name,
		Birthdate: author.Birthdate,
		CreatedAt: author.CreatedAt,
		UpdatedAt: author.UpdatedAt,
	})
}

// GetAuthors implements service.AuthorSvc.
func (svc *authorSvc) GetAuthors(ctx context.Context) *views.Response {
	author, err := svc.repo.GetAuthors(ctx)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	authors := make([]views.Author, 0)
	for _, ath := range author {
		authors = append(authors, views.Author{
			Id:        ath.Id,
			UserId:    ath.UserId,
			Name:      ath.Name,
			Birthdate: ath.Birthdate,
			CreatedAt: ath.CreatedAt,
			UpdatedAt: ath.UpdatedAt,
		})
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, authors)
}

// UpdateAuthor implements service.AuthorSvc.
func (svc *authorSvc) UpdateAuthor(ctx context.Context, author *params.UpdateAuthors, id uuid.UUID) *views.Response {
	a, err := svc.repo.GetAuthorById(ctx, id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	a.Name = author.Name
	a.Birthdate = author.Birthdate

	err = svc.repo.UpdateAuthor(ctx, a, id)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.UpdateAuthor{
		Id:        a.Id,
		UserId:    a.UserId,
		Name:      a.Name,
		Birthdate: a.Birthdate,
	})
}

func NewAuthorSvc(repo repository.AuthorRepo) service.AuthorSvc {
	return &authorSvc{
		repo: repo,
	}
}
