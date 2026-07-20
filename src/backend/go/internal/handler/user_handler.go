package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/security"
	"github.com/todo/backend/go/internal/service"
)

type UserHandler struct {
	svc service.UserService
	cfg *config.Config
}

func NewUserHandler(svc service.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{svc: svc, cfg: cfg}
}

func userServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrUsernameConflict), errors.Is(err, service.ErrEmailConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrInvalidCurrentPassword), errors.Is(err, service.ErrInactiveUser), errors.Is(err, service.ErrInvalidPasswordResetToken):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
func bind(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
func (h *UserHandler) currentUserID(c *gin.Context) (uint, bool) {
	id, err := security.AuthenticatedUserID(c.GetHeader("Authorization"), h.cfg.JWTSecretKey, time.Now().UTC())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return 0, false
	}
	return id, true
}
func (h *UserHandler) GetAll(c *gin.Context) {
	p := queryInt(c, "page", 1)
	s := queryInt(c, "pageSize", h.cfg.DefaultPageSize)
	r, e := h.svc.GetAll(p, s)
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) GetByID(c *gin.Context) {
	id, e := parseID(c)
	if e != nil {
		return
	}
	r, e := h.svc.GetByID(id)
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) Create(c *gin.Context) {
	var q dto.CreateUserRequest
	if !bind(c, &q) {
		return
	}
	r, e := h.svc.Create(q, auditUserID(c))
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusCreated, r)
}
func (h *UserHandler) Update(c *gin.Context) {
	id, e := parseID(c)
	if e != nil {
		return
	}
	var q dto.UpdateUserRequest
	if !bind(c, &q) {
		return
	}
	r, e := h.svc.Update(id, q, auditUserID(c))
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) Activate(c *gin.Context)   { h.setActive(c, true) }
func (h *UserHandler) Deactivate(c *gin.Context) { h.setActive(c, false) }
func (h *UserHandler) setActive(c *gin.Context, a bool) {
	id, e := parseID(c)
	if e != nil {
		return
	}
	r, e := h.svc.SetActive(id, a, auditUserID(c))
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) SignUp(c *gin.Context) {
	var q dto.SignUpRequest
	if !bind(c, &q) {
		return
	}
	r, e := h.svc.SignUp(q)
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusCreated, r)
}
func (h *UserHandler) GetProfile(c *gin.Context) {
	id, ok := h.currentUserID(c)
	if !ok {
		return
	}
	r, e := h.svc.GetProfile(id)
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, ok := h.currentUserID(c)
	if !ok {
		return
	}
	var q dto.UpdateProfileRequest
	if !bind(c, &q) {
		return
	}
	r, e := h.svc.UpdateProfile(id, q)
	if e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusOK, r)
}
func (h *UserHandler) ChangePassword(c *gin.Context) {
	id, ok := h.currentUserID(c)
	if !ok {
		return
	}
	var q dto.ChangePasswordRequest
	if !bind(c, &q) {
		return
	}
	if e := h.svc.ChangePassword(id, q); e != nil {
		userServiceError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}
func (h *UserHandler) RequestPasswordReset(c *gin.Context) {
	var q dto.ResetPasswordRequest
	if !bind(c, &q) {
		return
	}
	if e := h.svc.RequestPasswordReset(q); e != nil {
		userServiceError(c, e)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"message": "If the account exists, a password reset email has been queued."})
}
func (h *UserHandler) ConfirmPasswordReset(c *gin.Context) {
	var q dto.ConfirmPasswordResetRequest
	if !bind(c, &q) {
		return
	}
	if e := h.svc.ConfirmPasswordReset(q); e != nil {
		userServiceError(c, e)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) CreateToken(c *gin.Context) {
	var q dto.TokenRequest
	if !bind(c, &q) {
		return
	}
	r, err := h.svc.CreateToken(q)
	if errors.Is(err, service.ErrInvalidCredentials) {
		c.Header("WWW-Authenticate", "Bearer")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password."})
		return
	}
	if err != nil {
		userServiceError(c, err)
		return
	}
	c.Header("Cache-Control", "no-store")
	c.Header("Pragma", "no-cache")
	c.JSON(http.StatusOK, r)
}
