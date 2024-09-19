package gorm

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"gorm.io/gorm"
)

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) repository.UserRepo {
	return &userRepo{db: db}
}

// CreateUser implements repository.UserRepo.
func (repo *userRepo) CreateUser(ctx context.Context, user *models.User) error {
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}
	user.CreatedAt = time.Now()
	return repo.db.WithContext(ctx).Create(user).Error
}

// GetUserById implements repository.UserRepo.
func (repo *userRepo) GetUserById(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := new(models.User)
	return user, repo.db.WithContext(ctx).Where("id = ?", id).Take(user).Error
}

// GetUserById implements repository.UserRepo.
func (repo *userRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user := new(models.User)
	return user, repo.db.WithContext(ctx).Where("LOWER(username) = ?", strings.ToLower(username)).Take(user).Error
}

// GetUsers implements repository.UserRepo.
func (repo *userRepo) GetUsers(ctx context.Context, criteria map[string]interface{}) ([]*models.User, error) {
	var users []*models.User
	query := repo.db.WithContext(ctx).Model(&models.User{})

	for key, value := range criteria {
		query = query.Where(key+" = ?", value)
	}

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// UpdateUserById implements repository.UserRepo.
func (repo *userRepo) UpdateUserById(ctx context.Context, id uuid.UUID) error {
	user := new(models.User)
	user.UpdatedAt = time.Now()
	return repo.db.WithContext(ctx).Model(user).Where("id = ?", id).Updates(*user).Error
}
