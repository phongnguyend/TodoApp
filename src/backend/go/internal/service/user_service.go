package service

import (
	"errors"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/security"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound              = errors.New("user not found")
	ErrUsernameConflict          = errors.New("username is already in use")
	ErrEmailConflict             = errors.New("email is already in use")
	ErrInvalidCurrentPassword    = errors.New("the current password is incorrect")
	ErrInactiveUser              = errors.New("the user account is inactive")
	ErrInvalidPasswordResetToken = errors.New("the password reset token is invalid or expired")
	ErrInvalidCredentials        = errors.New("invalid email or password")
)

type UserService interface {
	GetAll(page, pageSize int) (dto.PaginatedResponse[dto.UserResponse], error)
	GetByID(id uint) (dto.UserResponse, error)
	Create(req dto.CreateUserRequest, actorUserID ...*uint) (dto.UserResponse, error)
	Update(id uint, req dto.UpdateUserRequest, actorUserID ...*uint) (dto.UserResponse, error)
	SetActive(id uint, active bool, actorUserID ...*uint) (dto.UserResponse, error)
	SignUp(req dto.SignUpRequest) (dto.UserResponse, error)
	GetProfile(id uint) (dto.UserResponse, error)
	UpdateProfile(id uint, req dto.UpdateProfileRequest) (dto.UserResponse, error)
	ChangePassword(id uint, req dto.ChangePasswordRequest) error
	RequestPasswordReset(req dto.ResetPasswordRequest) error
	ConfirmPasswordReset(req dto.ConfirmPasswordResetRequest) error
	CreateToken(req dto.TokenRequest) (dto.TokenResponse, error)
}

type userService struct {
	repo              repository.UserRepository
	cfg               *config.Config
	dummyPasswordHash string
}

func NewUserService(repo repository.UserRepository, cfg *config.Config) UserService {
	dummyHash, _ := security.HashPassword("not-a-real-password", cfg.PasswordHashIterations)
	return &userService{repo: repo, cfg: cfg, dummyPasswordHash: dummyHash}
}

func userResponse(u *models.User) dto.UserResponse {
	return dto.UserResponse{ID: u.ID, Username: u.Username, Email: u.Email, IsActive: u.IsActive,
		CreatedAt: u.CreatedAt, CreatedByUserID: u.CreatedByUserID,
		UpdatedAt: u.UpdatedAt, UpdatedByUserID: u.UpdatedByUserID}
}
func (s *userService) get(id uint) (*models.User, error) {
	u, err := s.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return u, err
}
func (s *userService) unique(username, email string, excluding *uint) error {
	exists, err := s.repo.UsernameExists(username, excluding)
	if err != nil {
		return err
	}
	if exists {
		return ErrUsernameConflict
	}
	exists, err = s.repo.EmailExists(email, excluding)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailConflict
	}
	return nil
}
func (s *userService) GetAll(page, size int) (dto.PaginatedResponse[dto.UserResponse], error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = s.cfg.DefaultPageSize
		if size < 1 {
			size = 20
		}
	}
	if s.cfg.MaxPageSize > 0 && size > s.cfg.MaxPageSize {
		size = s.cfg.MaxPageSize
	}
	r, err := s.repo.FindAll((page-1)*size, size)
	if err != nil {
		return dto.PaginatedResponse[dto.UserResponse]{}, err
	}
	items := make([]dto.UserResponse, len(r.Items))
	for i := range r.Items {
		items[i] = userResponse(&r.Items[i])
	}
	return dto.PaginatedResponse[dto.UserResponse]{Items: items, Total: r.Total, Page: page, PageSize: size, TotalPages: int(math.Ceil(float64(r.Total) / float64(size)))}, nil
}
func (s *userService) GetByID(id uint) (dto.UserResponse, error) {
	u, e := s.get(id)
	if e != nil {
		return dto.UserResponse{}, e
	}
	return userResponse(u), nil
}
func (s *userService) Create(req dto.CreateUserRequest, actorUserID ...*uint) (dto.UserResponse, error) {
	username, email := strings.TrimSpace(req.Username), strings.ToLower(strings.TrimSpace(req.Email))
	if err := s.unique(username, email, nil); err != nil {
		return dto.UserResponse{}, err
	}
	hash, err := security.HashPassword(req.Password, s.cfg.PasswordHashIterations)
	if err != nil {
		return dto.UserResponse{}, err
	}
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	u, err := s.repo.Create(&models.User{Username: username, Email: email, PasswordHash: hash,
		IsActive: active, CreatedByUserID: firstActor(actorUserID)})
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userResponse(u), nil
}
func (s *userService) Update(id uint, req dto.UpdateUserRequest, actorUserID ...*uint) (dto.UserResponse, error) {
	u, err := s.get(id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	username, email := u.Username, u.Email
	if req.Username != nil {
		username = strings.TrimSpace(*req.Username)
	}
	if req.Email != nil {
		email = strings.ToLower(strings.TrimSpace(*req.Email))
	}
	if err = s.unique(username, email, &id); err != nil {
		return dto.UserResponse{}, err
	}
	u.Username, u.Email = username, email
	if req.Password != nil {
		u.PasswordHash, err = security.HashPassword(*req.Password, s.cfg.PasswordHashIterations)
		if err != nil {
			return dto.UserResponse{}, err
		}
	}
	now := time.Now().UTC()
	u.UpdatedAt = &now
	u.UpdatedByUserID = firstActor(actorUserID)
	u, err = s.repo.Update(u)
	if err != nil {
		return dto.UserResponse{}, err
	}
	return userResponse(u), nil
}
func (s *userService) SetActive(id uint, active bool, actorUserID ...*uint) (dto.UserResponse, error) {
	u, e := s.get(id)
	if e != nil {
		return dto.UserResponse{}, e
	}
	u.IsActive = active
	now := time.Now().UTC()
	u.UpdatedAt = &now
	u.UpdatedByUserID = firstActor(actorUserID)
	u, e = s.repo.Update(u)
	if e != nil {
		return dto.UserResponse{}, e
	}
	return userResponse(u), nil
}
func (s *userService) SignUp(req dto.SignUpRequest) (dto.UserResponse, error) {
	active := true
	return s.Create(dto.CreateUserRequest{Username: req.Username, Email: req.Email, Password: req.Password, IsActive: &active})
}
func (s *userService) GetProfile(id uint) (dto.UserResponse, error) { return s.GetByID(id) }
func (s *userService) UpdateProfile(id uint, req dto.UpdateProfileRequest) (dto.UserResponse, error) {
	return s.Update(id, dto.UpdateUserRequest{Username: req.Username, Email: req.Email}, &id)
}
func (s *userService) ChangePassword(id uint, req dto.ChangePasswordRequest) error {
	u, e := s.get(id)
	if e != nil {
		return e
	}
	if !u.IsActive {
		return ErrInactiveUser
	}
	if !security.VerifyPassword(req.CurrentPassword, u.PasswordHash) {
		return ErrInvalidCurrentPassword
	}
	u.PasswordHash, e = security.HashPassword(req.NewPassword, s.cfg.PasswordHashIterations)
	if e != nil {
		return e
	}
	now := time.Now().UTC()
	u.UpdatedAt = &now
	u.UpdatedByUserID = &id
	_, e = s.repo.Update(u)
	return e
}
func (s *userService) RequestPasswordReset(req dto.ResetPasswordRequest) error {
	u, e := s.repo.FindByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	if errors.Is(e, gorm.ErrRecordNotFound) || u != nil && !u.IsActive {
		return nil
	}
	if e != nil {
		return e
	}
	token, e := security.CreateResetToken(u.ID, u.PasswordHash, s.cfg.PasswordResetSecretKey, time.Now().UTC().Add(time.Duration(s.cfg.PasswordResetTokenLifetimeMinutes)*time.Minute))
	if e != nil {
		return e
	}
	sep := "?"
	if strings.Contains(s.cfg.PasswordResetConfirmationURL, "?") {
		sep = "&"
	}
	resetURL := s.cfg.PasswordResetConfirmationURL + sep + "token=" + url.QueryEscape(token)
	_, e = s.repo.AddEmailLog(&models.EmailLog{Recipient: u.Email, Subject: "Reset your Todo API password", Body: fmt.Sprintf("Use this link to reset your password: %s\n\nThis link expires in %d minutes.", resetURL, s.cfg.PasswordResetTokenLifetimeMinutes), Status: "pending"})
	return e
}
func (s *userService) ConfirmPasswordReset(req dto.ConfirmPasswordResetRequest) error {
	id, fingerprint, e := security.DecodeResetToken(req.Token, s.cfg.PasswordResetSecretKey, time.Now().UTC())
	if e != nil {
		return ErrInvalidPasswordResetToken
	}
	u, e := s.get(id)
	if e != nil || !u.IsActive || security.PasswordFingerprint(u.PasswordHash) != fingerprint {
		return ErrInvalidPasswordResetToken
	}
	u.PasswordHash, e = security.HashPassword(req.NewPassword, s.cfg.PasswordHashIterations)
	if e != nil {
		return e
	}
	now := time.Now().UTC()
	u.UpdatedAt = &now
	u.UpdatedByUserID = nil
	_, e = s.repo.Update(u)
	return e
}

func (s *userService) CreateToken(req dto.TokenRequest) (dto.TokenResponse, error) {
	u, err := s.repo.FindByEmail(strings.ToLower(strings.TrimSpace(req.Email)))
	hash := s.dummyPasswordHash
	if err == nil && u != nil {
		hash = u.PasswordHash
	}
	valid := security.VerifyPassword(req.Password, hash)
	if errors.Is(err, gorm.ErrRecordNotFound) || err == nil && (u == nil || !valid || !u.IsActive) {
		return dto.TokenResponse{}, ErrInvalidCredentials
	}
	if err != nil {
		return dto.TokenResponse{}, err
	}
	now := time.Now().UTC()
	expiresIn := s.cfg.JWTTokenLifetimeMinutes * 60
	if expiresIn < 60 {
		expiresIn = 60
	}
	token, err := security.CreateJWT(u.ID, s.cfg.JWTSecretKey, now, now.Add(time.Duration(expiresIn)*time.Second))
	if err != nil {
		return dto.TokenResponse{}, err
	}
	return dto.TokenResponse{AccessToken: token, TokenType: "Bearer", ExpiresIn: expiresIn}, nil
}
