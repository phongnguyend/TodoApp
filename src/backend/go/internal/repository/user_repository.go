package repository

import (
	"strings"

	"github.com/todo/backend/go/internal/models"
	"gorm.io/gorm"
)

type userRepository struct{ db *gorm.DB }

func NewUserRepository(db *gorm.DB) UserRepository { return &userRepository{db: db} }

func (r *userRepository) FindAll(skip, limit int) (UserPaginatedResult, error) {
	var users []models.User
	var total int64
	if err := r.db.Model(&models.User{}).Count(&total).Error; err != nil {
		return UserPaginatedResult{}, err
	}
	if err := r.db.Offset(skip).Limit(limit).Order("created_at desc").Find(&users).Error; err != nil {
		return UserPaginatedResult{}, err
	}
	return UserPaginatedResult{Items: users, Total: total}, nil
}

func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("lower(email) = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) exists(column, value string, excludingID *uint) (bool, error) {
	q := r.db.Model(&models.User{}).Where("lower("+column+") = ?", strings.ToLower(value))
	if excludingID != nil {
		q = q.Where("id <> ?", *excludingID)
	}
	var count int64
	err := q.Count(&count).Error
	return count > 0, err
}

func (r *userRepository) UsernameExists(v string, id *uint) (bool, error) {
	return r.exists("username", v, id)
}
func (r *userRepository) EmailExists(v string, id *uint) (bool, error) {
	return r.exists("email", v, id)
}
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := r.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userRepository) Update(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
func (r *userRepository) AddEmailLog(log *models.EmailLog) (*models.EmailLog, error) {
	if err := r.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}
