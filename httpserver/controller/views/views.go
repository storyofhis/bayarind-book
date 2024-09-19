package views

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Payload interface{} `json:"payload,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

const (
	M_BAD_REQUEST                 = "BAD_REQUEST"
	M_INVALID_CREDENTIALS         = "INVALID_CREDENTIALS"
	M_CREATED                     = "CREATED"
	M_OK                          = "OK"
	M_USERNAME_ALREADY_USED       = "USER_ALREADY_USED"
	M_INTERNAL_SERVER_ERROR       = "INTERNAL_SERVER_ERROR"
	M_AUTHOR_SUCCESSFULLY_DELETED = "AUTHOR_SUCCESSFULLY_DELETED"
	M_AUTHOR_NOT_FOUND            = "AUTHOR_NOT_FOUND"
)

func SuccessResponse(status int, message string, payload interface{}) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Payload: payload,
	}
}

func ErrorReponse(status int, message string, error error) *Response {
	return &Response{
		Status:  status,
		Message: message,
		Error:   error.Error(),
	}
}

func WriteJsonResponse(ctx *gin.Context, res *Response) {
	ctx.JSON(res.Status, res)
}
