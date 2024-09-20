package user_controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/storyofhis/books-management/httpserver/service"
)

type UserController struct {
	svc service.UserSvc
}

func NewUserController(svc service.UserSvc) *UserController {
	return &UserController{
		svc: svc,
	}
}

func (control *UserController) Register(ctx *gin.Context) {
	var req params.Register
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = validator.New().Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := control.svc.Register(ctx, &req)
	views.WriteJsonResponse(ctx, response)
}

func (control *UserController) Login(ctx *gin.Context) {
	var req params.Login
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	err = validator.New().Struct(req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := control.svc.Login(ctx, &req)
	views.WriteJsonResponse(ctx, response)
}
