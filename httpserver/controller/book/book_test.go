package book_controller_test

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
	book_controller "github.com/storyofhis/books-management/httpserver/controller/book"
	"github.com/storyofhis/books-management/httpserver/controller/mocks"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateBook_Success(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.POST("/books", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.CreateBook(ctx)
	})

	payload := params.CreateBook{
		Title:    "Test Book",
		Isbn:     "123-456",
		AuthorId: uuid.New(),
	}
	body, _ := json.Marshal(payload)

	response := views.SuccessResponse(http.StatusCreated, views.M_CREATED, views.Book{
		Id:        uuid.New(),
		UserId:    uuid.New(),
		AuthorId:  payload.AuthorId,
		Title:     payload.Title,
		Isbn:      payload.Isbn,
		CreatedAt: time.Now(),
	})

	mockBookSvc.On("CreateBook", mock.Anything, mock.AnythingOfType("*params.CreateBook"), mock.AnythingOfType("uuid.UUID")).
		Return(response)

	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusCreated, rec.Code)
	mockBookSvc.AssertExpectations(t)
}

func TestCreateBook_InvalidInput(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/books", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.CreateBook(ctx)
	})

	payload := params.CreateBook{
		Isbn:     "123-456",
		AuthorId: uuid.New(),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockBookSvc.AssertNotCalled(t, "CreateBook")
}

func TestCreateBook_NoToken(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/books", controller.CreateBook)

	payload := params.CreateBook{
		Title:    "Test Book",
		Isbn:     "123-456",
		AuthorId: uuid.New(),
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	mockBookSvc.AssertNotCalled(t, "CreateBook")
}

func TestGetBooks_Success(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/books", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.GetBooks(ctx)
	})

	expectedBooks := []views.Book{
		{Id: uuid.New(), Title: "Book 1", Isbn: "123-456"},
		{Id: uuid.New(), Title: "Book 2", Isbn: "789-012"},
	}
	response := views.SuccessResponse(http.StatusOK, views.M_OK, expectedBooks)
	mockBookSvc.On("GetBooks", mock.Anything).Return(response)
	req, _ := http.NewRequest(http.MethodGet, "/books", nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":[{"id":"`+expectedBooks[0].Id.String()+`","title":"Book 1","isbn":"123-456"},{"id":"`+expectedBooks[1].Id.String()+`","title":"Book 2","isbn":"789-012"}]}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestGetBooks_EmptyResponse(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/books", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.GetBooks(ctx)
	})

	response := views.SuccessResponse(http.StatusOK, views.M_OK, []views.Book{})
	mockBookSvc.On("GetBooks", mock.Anything).Return(response)
	req, _ := http.NewRequest(http.MethodGet, "/books", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":[]}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestGetBookById_Success(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.GetBookById(ctx)
	})

	bookId := uuid.New()
	expectedBook := views.Book{
		Id:       bookId,
		Title:    "Test Book",
		Isbn:     "123-456",
		AuthorId: uuid.New(),
	}
	response := views.SuccessResponse(http.StatusOK, views.M_OK, expectedBook)
	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(response)

	req, _ := http.NewRequest(http.MethodGet, "/books/"+bookId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":{"id":"`+bookId.String()+`","title":"Test Book","isbn":"123-456"}}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestGetBookById_InvalidID(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/books/:id", controller.GetBookById)

	req, _ := http.NewRequest(http.MethodGet, "/books/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid book ID format"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "GetBookById")
}

func TestGetBookById_NotFound(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.GET("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.GetBookById(ctx)
	})

	bookId := uuid.New()
	response := views.ErrorReponse(http.StatusNotFound, "Book not found", errors.New("Book Not found"))
	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(response)
	req, _ := http.NewRequest(http.MethodGet, "/books/"+bookId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"status":404,"message":"Book not found"}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestUpdateBook_Success(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.UpdateBook(ctx)
	})

	bookId := uuid.New()
	updatePayload := params.UpdateBook{
		Title: "Updated Book",
		Isbn:  "456-789",
	}
	body, _ := json.Marshal(updatePayload)
	existingBook := views.Book{
		Id:     bookId,
		UserId: uuid.New(),
		Title:  "Old Book",
		Isbn:   "123-456",
	}
	bookResponse := views.SuccessResponse(http.StatusOK, views.M_OK, existingBook)
	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	updateResponse := views.SuccessResponse(http.StatusOK, views.M_OK, existingBook)
	mockBookSvc.On("UpdateBook", mock.Anything, mock.AnythingOfType("*params.UpdateBook"), bookId).
		Return(updateResponse)

	req, _ := http.NewRequest(http.MethodPut, "/books/"+bookId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":{"id":"`+bookId.String()+`","title":"Updated Book","isbn":"456-789"}}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestUpdateBook_InvalidID(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/books/:id", controller.UpdateBook)
	req, _ := http.NewRequest(http.MethodPut, "/books/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid book ID format"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "GetBookById")
}

func TestUpdateBook_InvalidPayload(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.UpdateBook(ctx)
	})

	bookId := uuid.New()
	req, _ := http.NewRequest(http.MethodPut, "/books/"+bookId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	mockBookSvc.AssertNotCalled(t, "GetBookById")
}

func TestUpdateBook_NotFound(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.UpdateBook(ctx)
	})

	bookId := uuid.New()
	bookResponse := views.ErrorReponse(http.StatusNotFound, "Book not found", errors.New("Book not found"))
	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)

	updatePayload := params.UpdateBook{
		Title: "Updated Book",
		Isbn:  "456-789",
	}
	body, _ := json.Marshal(updatePayload)

	req, _ := http.NewRequest(http.MethodPut, "/books/"+bookId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"status":404,"message":"Book not found"}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestUpdateBook_Unauthorized(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.PUT("/books/:id", controller.UpdateBook)
	bookId := uuid.New()
	updatePayload := params.UpdateBook{
		Title: "Updated Book",
		Isbn:  "456-789",
	}
	body, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest(http.MethodPut, "/books/"+bookId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"Token doesn't exist"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "GetBookById")
}

func TestUpdateBook_Forbidden(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.UpdateBook(ctx)
	})

	bookId := uuid.New()
	existingUserId := uuid.New()

	updatePayload := params.UpdateBook{
		Title: "Updated Book",
		Isbn:  "456-789",
	}
	body, _ := json.Marshal(updatePayload)

	bookResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Book{
		Id:     bookId,
		UserId: existingUserId,
		Title:  "Old Book",
		Isbn:   "123-456",
	})

	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	req, _ := http.NewRequest(http.MethodPut, "/books/"+bookId.String(), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.JSONEq(t, `{"error":"You do not have permission to update this book"}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestDeleteBook_Success(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.DeleteBook(ctx)
	})

	bookId := uuid.New()
	bookResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Book{
		Id:     bookId,
		UserId: uuid.New(),
	})

	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	deleteResponse := views.SuccessResponse(http.StatusOK, views.M_OK, nil)
	mockBookSvc.On("DeleteBook", mock.Anything, bookId).Return(deleteResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/books/"+bookId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `{"status":200,"message":"OK","data":null}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestDeleteBook_InvalidID(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/books/:id", controller.DeleteBook)

	req, _ := http.NewRequest(http.MethodDelete, "/books/invalid-id", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.JSONEq(t, `{"error":"Invalid book ID format"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "GetBookById")
}

func TestDeleteBook_NotFound(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	router.DELETE("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)

		controller.DeleteBook(ctx)
	})

	bookId := uuid.New()
	bookResponse := views.ErrorReponse(http.StatusNotFound, "Book not found", errors.New("Book Not Found"))
	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/books/"+bookId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"status":404,"message":"Book not found"}`, rec.Body.String())
	mockBookSvc.AssertExpectations(t)
}

func TestDeleteBook_Unauthorized(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/books/:id", controller.DeleteBook)
	bookId := uuid.New()
	bookResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Book{
		Id:     bookId,
		UserId: uuid.New(),
	})

	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/books/"+bookId.String(), nil)

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusUnauthorized, rec.Code)
	assert.JSONEq(t, `{"error":"Token doesn't exist"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "DeleteBook")
}

func TestDeleteBook_Forbidden(t *testing.T) {
	mockBookSvc := new(mocks.MockBookSvc)
	controller := book_controller.NewBookController(mockBookSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.DELETE("/books/:id", func(ctx *gin.Context) {
		claims := &common.CustomClaims{Id: uuid.New()}
		ctx.Set("userData", claims)
		controller.DeleteBook(ctx)
	})

	bookId := uuid.New()
	existingUserId := uuid.New()

	bookResponse := views.SuccessResponse(http.StatusOK, views.M_OK, views.Book{
		Id:     bookId,
		UserId: existingUserId,
	})

	mockBookSvc.On("GetBookById", mock.Anything, bookId).Return(bookResponse)
	req, _ := http.NewRequest(http.MethodDelete, "/books/"+bookId.String(), nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusForbidden, rec.Code)
	assert.JSONEq(t, `{"error":"You do not have permission to update this author"}`, rec.Body.String())
	mockBookSvc.AssertNotCalled(t, "DeleteBook")
}
