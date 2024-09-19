package httpserver

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/storyofhis/books-management/common"
	"github.com/storyofhis/books-management/httpserver/controller"
)

type router struct {
	router *gin.Engine

	user   *controller.UserController
	author *controller.AuthorController
	book   *controller.BookController
}

func NewRouter(r *gin.Engine, user *controller.UserController, author *controller.AuthorController, book *controller.BookController) *router {
	return &router{
		router: r,
		user:   user,
		author: author,
		book:   book,
	}
}

func (r *router) Start(port string) {
	r.router.POST("/auth/register", r.user.Register)
	r.router.POST("/auth/login", r.user.Login)

	r.router.POST("/authors", r.verifyToken, r.author.CreateAuthor)
	r.router.GET("/authors", r.verifyToken, r.author.GetAuthors)
	r.router.GET("/authors/:id", r.verifyToken, r.author.GetAuthorById)
	r.router.PUT("/authors/:id", r.verifyToken, r.author.UpdateAuthor)
	r.router.DELETE("/authors/:id", r.verifyToken, r.author.DeleteAuthor)

	r.router.POST("/books", r.verifyToken, r.book.CreateBook)
	r.router.GET("/books", r.verifyToken, r.book.GetBooks)
	r.router.GET("/books/:id", r.verifyToken, r.book.GetBookById)
	r.router.PUT("/books/:id", r.verifyToken, r.book.UpdateBook)
	r.router.DELETE("books/:id", r.verifyToken, r.book.DeleteBook)
	r.router.Run(port)
}

func (r *router) verifyToken(ctx *gin.Context) {
	bearerToken := strings.Split(ctx.Request.Header.Get("Authorization"), "Bearer ")
	if len(bearerToken) != 2 {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "invalid bearer token",
		})
		return
	}
	claims, err := common.ValidateToken(bearerToken[1])
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.Set("userData", claims)
}
