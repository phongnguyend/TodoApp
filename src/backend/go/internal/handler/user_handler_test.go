package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/handler"
	"github.com/todo/backend/go/internal/router"
)

type userServiceStub struct {
	signUp func(dto.SignUpRequest) (dto.UserResponse, error)
	reset  func(dto.ResetPasswordRequest) error
}

func (s *userServiceStub) GetAll(int, int) (dto.PaginatedResponse[dto.UserResponse], error) {
	return dto.PaginatedResponse[dto.UserResponse]{}, nil
}
func (s *userServiceStub) GetByID(uint) (dto.UserResponse, error) { return dto.UserResponse{}, nil }
func (s *userServiceStub) Create(dto.CreateUserRequest) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
func (s *userServiceStub) Update(uint, dto.UpdateUserRequest) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
func (s *userServiceStub) SetActive(uint, bool) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
func (s *userServiceStub) SignUp(q dto.SignUpRequest) (dto.UserResponse, error) { return s.signUp(q) }
func (s *userServiceStub) GetProfile(uint) (dto.UserResponse, error)            { return dto.UserResponse{}, nil }
func (s *userServiceStub) UpdateProfile(uint, dto.UpdateProfileRequest) (dto.UserResponse, error) {
	return dto.UserResponse{}, nil
}
func (s *userServiceStub) ChangePassword(uint, dto.ChangePasswordRequest) error       { return nil }
func (s *userServiceStub) RequestPasswordReset(q dto.ResetPasswordRequest) error      { return s.reset(q) }
func (s *userServiceStub) ConfirmPasswordReset(dto.ConfirmPasswordResetRequest) error { return nil }

func userRouter(s *userServiceStub) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	router.RegisterUsers(r, handler.NewUserHandler(s, &config.Config{DefaultPageSize: 20, JWTSecretKey: "test"}))
	return r
}
func userRequest(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestUserHandlerSignupReturnsCreated(t *testing.T) {
	s := &userServiceStub{signUp: func(q dto.SignUpRequest) (dto.UserResponse, error) {
		assert.Equal(t, "alice@example.com", q.Email)
		return dto.UserResponse{ID: 1, Username: q.Username, Email: q.Email, IsActive: true}, nil
	}, reset: func(dto.ResetPasswordRequest) error { return nil }}
	w := userRequest(userRouter(s), http.MethodPost, "/api/users/signup", `{"username":"alice","email":"alice@example.com","password":"password123"}`)
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), `"isActive":true`)
}

func TestUserHandlerRejectsInvalidSignup(t *testing.T) {
	called := false
	s := &userServiceStub{signUp: func(dto.SignUpRequest) (dto.UserResponse, error) { called = true; return dto.UserResponse{}, nil }, reset: func(dto.ResetPasswordRequest) error { return nil }}
	w := userRequest(userRouter(s), http.MethodPost, "/api/users/signup", `{"username":"alice","email":"invalid","password":"short"}`)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.False(t, called)
}

func TestUserHandlerProfileRequiresBearerToken(t *testing.T) {
	s := &userServiceStub{signUp: func(dto.SignUpRequest) (dto.UserResponse, error) { return dto.UserResponse{}, nil }, reset: func(dto.ResetPasswordRequest) error { return nil }}
	w := userRequest(userRouter(s), http.MethodGet, "/api/users/profile", "")
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUserHandlerPasswordResetAlwaysReturnsAccepted(t *testing.T) {
	s := &userServiceStub{signUp: func(dto.SignUpRequest) (dto.UserResponse, error) { return dto.UserResponse{}, nil }, reset: func(q dto.ResetPasswordRequest) error { assert.Equal(t, "missing@example.com", q.Email); return nil }}
	w := userRequest(userRouter(s), http.MethodPost, "/api/users/password/reset", `{"email":"missing@example.com"}`)
	assert.Equal(t, http.StatusAccepted, w.Code)
	assert.Contains(t, w.Body.String(), "account exists")
}
