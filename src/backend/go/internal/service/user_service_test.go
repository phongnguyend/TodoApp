package service_test

import (
	"fmt"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/security"
	"github.com/todo/backend/go/internal/service"
	"gorm.io/gorm"
)

type userRepoFake struct {
	users map[uint]*models.User
	next  uint
	logs  []*models.EmailLog
}

func newUserRepoFake() *userRepoFake { return &userRepoFake{users: map[uint]*models.User{}, next: 1} }
func (r *userRepoFake) FindAll(skip, limit int) (repository.UserPaginatedResult, error) {
	all := []models.User{}
	for _, u := range r.users {
		all = append(all, *u)
	}
	end := skip + limit
	if end > len(all) {
		end = len(all)
	}
	if skip > len(all) {
		skip = len(all)
	}
	return repository.UserPaginatedResult{Items: all[skip:end], Total: int64(len(all))}, nil
}
func (r *userRepoFake) FindByID(id uint) (*models.User, error) {
	u, ok := r.users[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return u, nil
}
func (r *userRepoFake) FindByEmail(email string) (*models.User, error) {
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}
func (r *userRepoFake) UsernameExists(v string, x *uint) (bool, error) {
	for id, u := range r.users {
		if u.Username == v && (x == nil || id != *x) {
			return true, nil
		}
	}
	return false, nil
}
func (r *userRepoFake) EmailExists(v string, x *uint) (bool, error) {
	for id, u := range r.users {
		if u.Email == v && (x == nil || id != *x) {
			return true, nil
		}
	}
	return false, nil
}
func (r *userRepoFake) Create(u *models.User) (*models.User, error) {
	u.ID = r.next
	r.next++
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now().UTC()
	}
	r.users[u.ID] = u
	return u, nil
}
func (r *userRepoFake) Update(u *models.User) (*models.User, error) { r.users[u.ID] = u; return u, nil }
func (r *userRepoFake) AddEmailLog(l *models.EmailLog) (*models.EmailLog, error) {
	r.logs = append(r.logs, l)
	return l, nil
}

func testUserConfig() *config.Config {
	return &config.Config{DefaultPageSize: 20, MaxPageSize: 100, PasswordHashIterations: 1000, JWTSecretKey: "test", PasswordResetSecretKey: "reset-test", PasswordResetTokenLifetimeMinutes: 60, PasswordResetConfirmationURL: "/reset-password"}
}

func TestUserServiceCreateNormalizesAndHashesPassword(t *testing.T) {
	r := newUserRepoFake()
	s := service.NewUserService(r, testUserConfig())
	active := false
	got, err := s.Create(dto.CreateUserRequest{Username: " Alice ", Email: "ALICE@Example.com", Password: "password123", IsActive: &active})
	require.NoError(t, err)
	assert.Equal(t, "Alice", got.Username)
	assert.Equal(t, "alice@example.com", got.Email)
	assert.False(t, got.IsActive)
	assert.True(t, security.VerifyPassword("password123", r.users[got.ID].PasswordHash))
	assert.NotEqual(t, "password123", r.users[got.ID].PasswordHash)
}

func TestUserServiceRejectsDuplicateUsername(t *testing.T) {
	r := newUserRepoFake()
	s := service.NewUserService(r, testUserConfig())
	_, err := s.SignUp(dto.SignUpRequest{Username: "alice", Email: "alice@example.com", Password: "password123"})
	require.NoError(t, err)
	_, err = s.SignUp(dto.SignUpRequest{Username: "alice", Email: "other@example.com", Password: "password123"})
	assert.ErrorIs(t, err, service.ErrUsernameConflict)
}

func TestUserServiceChangeAndResetPassword(t *testing.T) {
	r := newUserRepoFake()
	s := service.NewUserService(r, testUserConfig())
	u, err := s.SignUp(dto.SignUpRequest{Username: "alice", Email: "alice@example.com", Password: "password123"})
	require.NoError(t, err)
	assert.ErrorIs(t, s.ChangePassword(u.ID, dto.ChangePasswordRequest{CurrentPassword: "wrong", NewPassword: "new-password123"}), service.ErrInvalidCurrentPassword)
	require.NoError(t, s.RequestPasswordReset(dto.ResetPasswordRequest{Email: u.Email}))
	require.Len(t, r.logs, 1)
	var token string
	_, err = fmt.Sscanf(r.logs[0].Body, "Use this link to reset your password: /reset-password?token=%s\n", &token)
	require.NoError(t, err)
	decoded, err := url.QueryUnescape(token)
	require.NoError(t, err)
	require.NoError(t, s.ConfirmPasswordReset(dto.ConfirmPasswordResetRequest{Token: decoded, NewPassword: "new-password123"}))
	assert.True(t, security.VerifyPassword("new-password123", r.users[u.ID].PasswordHash))
}

func TestUserServiceUnknownResetEmailDoesNotLeakAccountExistence(t *testing.T) {
	r := newUserRepoFake()
	s := service.NewUserService(r, testUserConfig())
	assert.NoError(t, s.RequestPasswordReset(dto.ResetPasswordRequest{Email: "missing@example.com"}))
	assert.Empty(t, r.logs)
}

var _ repository.UserRepository = (*userRepoFake)(nil)
