package controller

import "github.com/gin-gonic/gin"

type UserController interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthorController interface {
	CreateAuthor(ctx *gin.Context)
	GetAuthors(ctx *gin.Context)
	GetAuthorById(ctx *gin.Context)
	UpdateAuthor(ctx *gin.Context)
	DeleteAuthor(ctx *gin.Context)
}

type BookController interface {
	CreateBook(ctx *gin.Context)
	GetBooks(ctx *gin.Context)
	GetBookById(ctx *gin.Context)
	UpdateBook(ctx *gin.Context)
	DeleteBook(ctx *gin.Context)
}
