package user

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/storyofhis/books-management/common"
	"github.com/storyofhis/books-management/config"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"github.com/storyofhis/books-management/httpserver/service"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userSvc struct {
	repo repository.UserRepo
}

// Register implements service.UserSvc.
func (svc *userSvc) Register(ctx context.Context, user *params.Register) *views.Response {
	_, err := svc.repo.GetUserByUsername(ctx, user.Username)
	if err == nil {
		return views.ErrorReponse(http.StatusBadRequest, views.M_USERNAME_ALREADY_USED, err)
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	input := models.User{
		Username: user.Username,
		Password: string(hashed),
	}

	err = svc.repo.CreateUser(ctx, &input)
	if err != nil {
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Register{
		Id:        input.Id,
		Username:  input.Username,
		Password:  input.Password,
		CreatedAt: input.CreatedAt,
	})
}

// Login implements service.UserSvc.
func (svc *userSvc) Login(ctx context.Context, user *params.Login) *views.Response {
	model, err := svc.repo.GetUserByUsername(ctx, user.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorReponse(http.StatusBadRequest, views.M_INVALID_CREDENTIALS, err)
		}
		return views.ErrorReponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(user.Password))
	if err != nil {
		return views.ErrorReponse(http.StatusBadRequest, views.M_INVALID_CREDENTIALS, err)
	}

	claims := &common.CustomClaims{
		Id: model.Id,
	}
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(config.GetJwtExpiredTime())).Unix()
	claims.Subject = model.Username

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(config.GetJwtSignature())
	return views.SuccessResponse(http.StatusOK, views.M_OK, views.Login{
		Id:       model.Id,
		Username: model.Username,
		Token:    ss,
	})
}

func NewUserSvc(repo repository.UserRepo) service.UserSvc {
	return &userSvc{
		repo: repo,
	}
}
