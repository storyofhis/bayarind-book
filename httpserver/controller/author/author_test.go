package author_controller_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/storyofhis/books-management/common"
	author_controller "github.com/storyofhis/books-management/httpserver/controller/author"
	"github.com/storyofhis/books-management/httpserver/controller/mocks"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAuthor_Success(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/authors", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.CreateAuthor(ctx)
	})

	payload := params.CreateAuthors{
		Name:      "Test Author",
		Birthdate: time.Now(),
	}

	body, _ := json.Marshal(payload)
	response := views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Author{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		Name:      payload.Name,
		Birthdate: payload.Birthdate,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})

	mockAuthorSvc.On("CreateAuhtor", mock.Anything, mock.AnythingOfType("*params.CreateAuthor"), mock.AnythingOfType("uuid.UUID")).
		Return(response)
	req, _ := http.NewRequest(http.MethodPost, "/authors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockAuthorSvc.AssertExpectations(t)
}

func TestCreateAuthor_InvalidInput(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/authors", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.CreateAuthor(ctx)
	})

	payload := params.CreateAuthors{
		Birthdate: time.Now(),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/authors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockAuthorSvc.AssertNotCalled(t, "CreateAuthor")
}

func TestCreateAuthor_NoToken(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/authors", controller.CreateAuthor)

	payload := params.CreateAuthors{
		Name:      "Test Author",
		Birthdate: time.Now(),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/authors", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockAuthorSvc.AssertNotCalled(t, "CreateAuthor")
}

func TestGetAuthors_Success(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/authors", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.GetAuthors(ctx)
	})

	expectedAuthors := []views.Author{
		{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			Name:      "Author 1",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Id:        uuid.New(),
			UserId:    uuid.New(),
			Name:      "Author 2",
			Birthdate: time.Now(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	response := views.SuccessResponse(http.StatusOK, views.M_OK, expectedAuthors)
	mockAuthorSvc.On("GetAuthors", mock.Anything).Return(response)

	req, _ := http.NewRequest(http.MethodGet, "/authors", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponse, _ := json.Marshal(response)
	assert.JSONEq(t, string(expectedResponse), rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestGetAuthors_EmptyResponse(t *testing.T) {
	moockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(moockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/authors", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.GetAuthors(ctx)
	})

	response := views.SuccessResponse(http.StatusOK, views.M_OK, []views.Author{})
	moockAuthorSvc.On("GetAuthors", mock.Anything).Return(response)
	req, _ := http.NewRequest(http.MethodGet, "/authors", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	expectedResponse, _ := json.Marshal(response)
	assert.JSONEq(t, string(expectedResponse), rec.Body.String())
	moockAuthorSvc.AssertExpectations(t)
}

func TestGetAuthorById_Success(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.GetAuthorById(ctx)
	})

	authorId := uuid.New()
	expectedAuthor := views.Author{
		Id:        authorId,
		UserId:    uuid.New(),
		Name:      "Author 1",
		Birthdate: time.Now(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	response := views.SuccessResponse(http.StatusOK, views.M_OK, expectedAuthor)
	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(response)

	req, _ := http.NewRequest(http.MethodGet, "/authors/"+authorId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	mockAuthorSvc.AssertExpectations(t)
}

func TestGetAuthorById_InvalidId(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/authors/:id", controller.GetAuthorById)

	req, _ := http.NewRequest(http.MethodGet, "/authors/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockAuthorSvc.AssertNotCalled(t, "GetAuthorById")
}

func TestGetAuthorById_NotFound(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.GetAuthorById(ctx)
	})

	authorId := uuid.New()
	response := views.ErrorReponse(http.StatusNotFound, "Author not found", errors.New("Author not found"))
	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(response)
	req, _ := http.NewRequest(http.MethodGet, "/authors/"+authorId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockAuthorSvc.AssertExpectations(t)
}

func TestUpdateAuthor_Success(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{
			Id: uuid.New(),
		}
		ctx.Set("userData", claims)
		controller.UpdateAuthor(ctx)
	})

	authorId := uuid.New()
	updatePayload := params.UpdateAuthors{
		Name:      "update author",
		Birthdate: time.Now(),
		UpdateAt:  time.Now(),
	}
	body, _ := json.Marshal(updatePayload)
	existingAuthor := views.Author{
		Id:        authorId,
		UserId:    uuid.New(),
		Name:      "old author",
		Birthdate: time.Now(),
	}
	authorResponse := views.SuccessResponse(http.StatusOK, views.M_OK, existingAuthor)
	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)
	updateResponse := views.SuccessResponse(http.StatusOK, views.M_OK, existingAuthor)
	mockAuthorSvc.On("UpdateAuthor", mock.Anything, mock.AnythingOfType("*params.UpdateAuthor")).
		Return(updateResponse)
	req, _ := http.NewRequest(http.MethodPut, "/authors/"+authorId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponse, _ := json.Marshal(updateResponse)

	// Compare the actual response with the expected JSON response
	assert.JSONEq(t, string(expectedResponse), rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestUpdateAuthor_InvalidID(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/authors/:id", controller.UpdateAuthor)
	req, _ := http.NewRequest(http.MethodPut, "/authors/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	assert.JSONEq(t, `{"error":"Invalid author ID format"}`, rec.Body.String())
	mockAuthorSvc.AssertNotCalled(t, "GetAuthorById")
}

func TestUpdateAuthor_InvalidPayload(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.UpdateAuthor(ctx)
	})

	authorId := uuid.New()
	req, _ := http.NewRequest(http.MethodPut, "/authors/"+authorId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockAuthorSvc.AssertNotCalled(t, "GetAuthorkById")
}

func TestUpdateAuthor_NotFound(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.UpdateAuthor(ctx)
	})

	authorId := uuid.New()
	authorResponse := views.ErrorReponse(http.StatusNotFound, "Author not found", errors.New("Author not found"))
	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)

	updatePayload := params.UpdateAuthors{
		Name:      "Updated Author",
		Birthdate: time.Now(),
		UpdateAt:  time.Now(),
	}
	body, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, "/authors/"+authorId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"status":404,"message":"Author not found"}`, rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestUpdateAuthor_Unauthorized(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/authors/:id", controller.UpdateAuthor)
	authorId := uuid.New()
	updatePayload := params.UpdateAuthors{
		Name:      "Updated Author",
		Birthdate: time.Now(),
	}
	body, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/authors/"+authorId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"Token doesn't exist"}`, rec.Body.String())
	mockAuthorSvc.AssertNotCalled(t, "GetAuthorById")
}

func TestUpdateAuthor_Forbidden(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.UpdateAuthor(ctx)
	})

	authorId := uuid.New()
	existingUserId := uuid.New()

	updatePayload := params.UpdateAuthors{
		Name:      "Updated Author",
		Birthdate: time.Now(),
	}
	body, _ := json.Marshal(updatePayload)

	authorResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Author{
		Id:        authorId,
		UserId:    existingUserId,
		Name:      "Old Author",
		Birthdate: time.Now(),
	})

	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)
	req, _ := http.NewRequest(http.MethodPut, "/authors/"+authorId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.JSONEq(t, `{"error":"You do not have permission to update this author"}`, rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestDeleteAuthor_Success(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.DeleteAuthor(ctx)
	})

	authorId := uuid.New()
	authorResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Author{
		Id:     authorId,
		UserId: uuid.New(),
	})

	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)
	deleteResponse := views.SuccessResponse(http.StatusOK, views.M_OK, nil)
	mockAuthorSvc.On("DeleteAuthor", mock.Anything, authorId).Return(deleteResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/authors/"+authorId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":null}`, rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestDeleteAuthor_InvalidID(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/authors/:id", controller.DeleteAuthor)

	req, _ := http.NewRequest(http.MethodDelete, "/authors/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid author ID format"}`, rec.Body.String())
	mockAuthorSvc.AssertNotCalled(t, "GetAuthorById")
}

func TestDeleteAuthor_NotFound(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.DeleteAuthor(ctx)
	})

	authorId := uuid.New()
	authorResponse := views.ErrorReponse(http.StatusNotFound, "Author not found", errors.New("Author Not Found"))
	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/authors/"+authorId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"status":404,"message":"Author not found"}`, rec.Body.String())
	mockAuthorSvc.AssertExpectations(t)
}

func TestDeleteAuthor_Unauthorized(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/authors/:id", controller.DeleteAuthor)
	authorId := uuid.New()
	authoresponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Author{
		Id:     authorId,
		UserId: uuid.New(),
	})

	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authoresponse)
	req, _ := http.NewRequest(http.MethodDelete, "/authors/"+authorId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"Token doesn't exist"}`, rec.Body.String())
	mockAuthorSvc.AssertNotCalled(t, "DeleteAuthor")
}

func TestDeleteAuthor_Forbidden(t *testing.T) {
	mockAuthorSvc := new(mocks.MockAuthorSvc)
	controller := author_controller.NewAuthorController(mockAuthorSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/authors/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.DeleteAuthor(ctx)
	})

	authorId := uuid.New()
	existingUserId := uuid.New()

	authorResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Author{
		Id:     authorId,
		UserId: existingUserId,
	})

	mockAuthorSvc.On("GetAuthorById", mock.Anything, authorId).Return(authorResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/authors/"+authorId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.JSONEq(t, `{"error":"You do not have permission to update this author"}`, rec.Body.String())
	mockAuthorSvc.AssertNotCalled(t, "DeleteAuthor")
}
