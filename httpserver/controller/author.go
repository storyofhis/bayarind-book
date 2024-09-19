package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/storyofhis/books-management/common"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/service"
)

type AuthorController struct {
	svc service.AuthorSvc
}

func NewAuthorController(svc service.AuthorSvc) *AuthorController {
	return &AuthorController{
		svc: svc,
	}
}

func (control *AuthorController) CreateAuthor(ctx *gin.Context) {
	var req params.CreateAuthors
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	claims, exist := ctx.Get("userData")
	if !exist {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "token doesn't exist",
		})
		return
	}

	userData := claims.(*common.CustomClaims)
	userId := userData.Id
	err := validator.New().Struct(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := control.svc.CreateAuthor(ctx, &req, userId)
	views.WriteJsonResponse(ctx, response)
}

func (control *AuthorController) GetAuthors(ctx *gin.Context) {
	reponse := control.svc.GetAuthors(ctx)
	views.WriteJsonResponse(ctx, reponse)
}

func (control *AuthorController) GetAuthorById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	authorId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid author ID format",
		})
		return
	}

	authorResponse := control.svc.GetAuthorById(ctx, authorId)
	if authorResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, authorResponse)
		return
	}

	views.WriteJsonResponse(ctx, authorResponse)
}

func (control *AuthorController) UpdateAuthor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	authorId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid author ID format",
		})
		return
	}

	var req params.UpdateAuthors
	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	authorResponse := control.svc.GetAuthorById(ctx, authorId)
	if authorResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, authorResponse)
		return
	}

	authorDetails, ok := authorResponse.Payload.(views.Author)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process author details",
		})
		return
	}

	claims, exists := ctx.Get("userData")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token doesn't exist",
		})
		return
	}

	userData := claims.(*common.CustomClaims)
	if userData.Id != authorDetails.UserId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update this author",
		})
		return
	}
	response := control.svc.UpdateAuthor(ctx, &req, authorId)
	views.WriteJsonResponse(ctx, response)
}

func (control *AuthorController) DeleteAuthor(ctx *gin.Context) {
	idParam := ctx.Param("id")
	authorId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid author ID format",
		})
		return
	}

	authorResponse := control.svc.GetAuthorById(ctx, authorId)
	if authorResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, authorResponse)
		return
	}

	authorDetails, ok := authorResponse.Payload.(views.Author)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process author details",
		})
		return
	}

	claims, exists := ctx.Get("userData")
	if !exists {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "Token doesn't exist",
		})
		return
	}

	userData := claims.(*common.CustomClaims)

	if userData.Id != authorDetails.UserId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update this author",
		})
		return
	}

	reponse := control.svc.DeleteAuthor(ctx, authorId)
	views.WriteJsonResponse(ctx, reponse)
}
