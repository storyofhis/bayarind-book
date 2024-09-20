package book_controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/storyofhis/books-management/common"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/service"
)

type BookController struct {
	svc service.BookSvc
}

func NewBookController(svc service.BookSvc) *BookController {
	return &BookController{
		svc: svc,
	}
}

func (control *BookController) CreateBook(ctx *gin.Context) {
	var req params.CreateBook
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

	reponse := control.svc.CreateBook(ctx, &req, userId)
	views.WriteJsonResponse(ctx, reponse)
}

func (control *BookController) GetBooks(ctx *gin.Context) {
	reponse := control.svc.GetBooks(ctx)
	views.WriteJsonResponse(ctx, reponse)
}

func (control *BookController) GetBookById(ctx *gin.Context) {
	idParam := ctx.Param("id")
	bookId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID format",
		})
		return
	}

	bookResponse := control.svc.GetBookById(ctx, bookId)
	if bookResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, bookResponse)
		return
	}

	views.WriteJsonResponse(ctx, bookResponse)
}

func (control *BookController) UpdateBook(ctx *gin.Context) {
	idParam := ctx.Param("id")
	bookId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID format",
		})
		return
	}

	var req params.UpdateBook
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	bookResponse := control.svc.GetBookById(ctx, bookId)
	if bookResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, bookResponse)
		return
	}

	bookDetails, ok := bookResponse.Payload.(views.Book)
	fmt.Println("user_id : ", bookDetails.UserId)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process book details",
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
	if userData.Id != bookDetails.UserId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update this book",
		})
		return
	}
	response := control.svc.UpdateBook(ctx, &req, bookId)
	views.WriteJsonResponse(ctx, response)
}

func (control *BookController) DeleteBook(ctx *gin.Context) {
	idParam := ctx.Param("id")
	bookId, err := uuid.Parse(idParam)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Invalid book ID format",
		})
		return
	}

	bookResponse := control.svc.GetBookById(ctx, bookId)
	if bookResponse.Status != http.StatusOK {
		views.WriteJsonResponse(ctx, bookResponse)
		return
	}

	bookDetails, ok := bookResponse.Payload.(views.Book)
	if !ok {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to process book details",
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

	if userData.Id != bookDetails.UserId {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"error": "You do not have permission to update this author",
		})
		return
	}

	response := control.svc.DeleteBook(ctx, bookId)
	views.WriteJsonResponse(ctx, response)
}
