package main

import (
	"github.com/gin-gonic/gin"
	"github.com/storyofhis/books-management/config"
	"github.com/storyofhis/books-management/httpserver"
	"github.com/storyofhis/books-management/httpserver/controller"
	"github.com/storyofhis/books-management/httpserver/repository/gorm"
	"github.com/storyofhis/books-management/httpserver/service/author"
	"github.com/storyofhis/books-management/httpserver/service/book"
	"github.com/storyofhis/books-management/httpserver/service/user"
)

func main() {
	db, err := config.ConnectGorm()
	if err != nil {
		panic(err)
	}

	router := gin.Default()
	config.GenerateJwtSignature()

	userRepo := gorm.NewUserRepo(db)
	userSvc := user.NewUserSvc(userRepo)
	userControl := controller.NewUserController(userSvc)

	authorRepo := gorm.NewAuthorRepo(db)
	authorSvc := author.NewAuthorSvc(authorRepo)
	authorControl := controller.NewAuthorController(authorSvc)

	bookRepo := gorm.NewBookRepo(db)
	bookSvc := book.NewBookSvc(bookRepo)
	bookControl := controller.NewBookController(bookSvc)

	app := httpserver.NewRouter(router, userControl, authorControl, bookControl)
	app.Start(":" + "8080")
}
