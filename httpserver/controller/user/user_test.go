package user_controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/storyofhis/books-management/httpserver/controller/mocks"
	"github.com/storyofhis/books-management/httpserver/controller/params"
	user_controller "github.com/storyofhis/books-management/httpserver/controller/user"
	"github.com/storyofhis/books-management/httpserver/controller/views"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister_Success(t *testing.T) {
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/register", controller.Register)
	registerPayload := params.Register{
		Username: "instagram",
		Password: "password123",
	}
	body, _ := json.Marshal(registerPayload)
	expectedResponse := views.SuccessResponse(http.StatusOK, views.M_OK, "User registered successfully")
	mockUserSvc.On("Register", mock.Anything, mock.AnythingOfType("*params.Register")).Return(expectedResponse)

	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseJSON), rec.Body.String())
	mockUserSvc.AssertExpectations(t)
}

func TestRegister_ShouldBindJSONError(t *testing.T) {
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/register", controller.Register)
	invalidJSON := []byte(`{ "Username": "instagram", "Password": `)

	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expectedResponse := `{"error":"unexpected end of JSON input"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	mockUserSvc.AssertNotCalled(t, "Register", mock.Anything, mock.AnythingOfType("*params.Register"))
}

func TestRegister_ValidationError(t *testing.T) {
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/register", controller.Register)

	invalidPayload := params.Register{
		Username: "instagram",
	}
	body, _ := json.Marshal(invalidPayload)

	req, _ := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	expectedResponse := `{"error":"Key: 'Register.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())
	mockUserSvc.AssertNotCalled(t, "Register", mock.Anything, mock.AnythingOfType("*params.Register"))
}

func TestLogin(t *testing.T) {
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/login", controller.Login)

	loginPayload := params.Login{
		Username: "instagram",
		Password: "password123",
	}
	body, _ := json.Marshal(loginPayload)

	expectedResponse := views.SuccessResponse(http.StatusOK, views.M_OK, "User logged in successfully")
	mockUserSvc.On("Login", mock.Anything, mock.AnythingOfType("*params.Login")).Return(expectedResponse)

	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	expectedResponseJSON, _ := json.Marshal(expectedResponse)
	assert.JSONEq(t, string(expectedResponseJSON), rec.Body.String())

	mockUserSvc.AssertExpectations(t)
}

func TestLogin_ShouldBindJSONError(t *testing.T) {
	// Mock the service
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)

	// Set up Gin and router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/login", controller.Login)

	// Create an invalid request payload (invalid JSON)
	invalidJSON := []byte(`{ "Username": "instagram", "Password": `) // Missing closing quotes and value for password

	// Create the HTTP POST request with invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert the response code
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Assert the response contains the JSON binding error
	expectedResponse := `{"error":"unexpected end of JSON input"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	// Verify that the mock service's Login method was not called
	mockUserSvc.AssertNotCalled(t, "Login", mock.Anything, mock.AnythingOfType("*params.Login"))
}

func TestLogin_ValidationError(t *testing.T) {
	// Mock the service
	mockUserSvc := new(mocks.MockUserSvc)
	controller := user_controller.NewUserController(mockUserSvc)

	// Set up Gin and router
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/login", controller.Login)

	// Create an invalid request payload (missing Password field)
	invalidPayload := params.Login{
		Username: "instagram", // Missing Password field
	}
	body, _ := json.Marshal(invalidPayload)

	// Create the HTTP POST request
	req, _ := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Assert the response code
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Assert the response contains the validation error
	expectedResponse := `{"error":"Key: 'Login.Password' Error:Field validation for 'Password' failed on the 'required' tag"}`
	assert.JSONEq(t, expectedResponse, rec.Body.String())

	// Verify that the mock service's Login method was not called
	mockUserSvc.AssertNotCalled(t, "Login", mock.Anything, mock.AnythingOfType("*params.Login"))
}
