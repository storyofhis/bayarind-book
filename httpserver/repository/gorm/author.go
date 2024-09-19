package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"gorm.io/gorm"
)

type authorRepo struct {
	db *gorm.DB
}

func NewAuthorRepo(db *gorm.DB) repository.AuthorRepo {
	return &authorRepo{db: db}
}

// CreateAuthor implements repository.AuthorRepo.
func (repo *authorRepo) CreateAuthor(ctx context.Context, author *models.Author) error {
	author.Id = uuid.New()
	author.CreatedAt = time.Now()
	return repo.db.WithContext(ctx).Create(author).Error
}

// DeleteAuthor implements repository.AuthorRepo.
func (repo *authorRepo) DeleteAuthor(ctx context.Context, id uuid.UUID) error {
	return repo.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Author{}).Error
}

// GetAuthorById implements repository.AuthorRepo.
func (repo *authorRepo) GetAuthorById(ctx context.Context, id uuid.UUID) (*models.Author, error) {
	author := new(models.Author)
	return author, repo.db.WithContext(ctx).Where("id = ?", id).Take(author).Error
}

// GetAuthors implements repository.AuthorRepo.
func (repo *authorRepo) GetAuthors(ctx context.Context) ([]*models.Author, error) {
	var authors []*models.Author

	err := repo.db.WithContext(ctx).Find(&authors).Error
	if err != nil {
		return nil, err
	}
	return authors, nil
}

// UpdateAuthor implements repository.AuthorRepo.
func (repo *authorRepo) UpdateAuthor(ctx context.Context, author *models.Author, id uuid.UUID) error {
	author.UpdatedAt = time.Now()
	return repo.db.WithContext(ctx).Model(author).Where("id = ?", id).Updates(author).Error
}
