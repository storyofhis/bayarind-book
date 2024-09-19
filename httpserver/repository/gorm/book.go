package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/storyofhis/books-management/httpserver/repository"
	"github.com/storyofhis/books-management/httpserver/repository/models"
	"gorm.io/gorm"
)

type bookRepo struct {
	db *gorm.DB
}

func NewBookRepo(db *gorm.DB) repository.BookRepo {
	return &bookRepo{
		db: db,
	}
}

// CreateBook implements repository.BookRepo.
func (repo *bookRepo) CreateBook(ctx context.Context, book *models.Book) error {
	book.Id = uuid.New()
	book.CreatedAt = time.Now()
	return repo.db.WithContext(ctx).Create(book).Error
}

// DeleteBook implements repository.BookRepo.
func (repo *bookRepo) DeleteBook(ctx context.Context, id uuid.UUID) error {
	return repo.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Book{}).Error
}

// GetBookById implements repository.BookRepo.
func (repo *bookRepo) GetBookById(ctx context.Context, id uuid.UUID) (*models.Book, error) {
	book := new(models.Book)
	return book, repo.db.WithContext(ctx).Where("id = ?", id).Take(book).Error
}

// GetBooks implements repository.BookRepo.
func (repo *bookRepo) GetBooks(ctx context.Context) ([]*models.Book, error) {
	var books []*models.Book

	err := repo.db.WithContext(ctx).Find(&books).Error
	if err != nil {
		return nil, err
	}
	return books, nil
}

// UpdateBook implements repository.BookRepo.
func (repo *bookRepo) UpdateBook(ctx context.Context, book *models.Book, id uuid.UUID) error {
	book.UpdatedAt = time.Now()
	return repo.db.WithContext(ctx).Model(book).Where("id = ?", id).Updates(book).Error
}
