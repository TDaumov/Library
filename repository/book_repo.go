package repository

import (
	"context"
	"library_vebservice/models"

	"gorm.io/gorm"
)

type BookRepo interface {
	Create(ctx context.Context, book *models.Book) error
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Book, error)
	List(ctx context.Context) ([]models.Book, error)
	CreateBorrow(ctx context.Context, borrow *models.Borrow) error // ← добавляем
}

type bookRepo struct{ db *gorm.DB }

func NewBookRepo(db *gorm.DB) BookRepo { return &bookRepo{db: db} }

func (r *bookRepo) Create(ctx context.Context, book *models.Book) error {
	return r.db.WithContext(ctx).Create(book).Error
}

func (r *bookRepo) Update(ctx context.Context, book *models.Book) error {
	return r.db.WithContext(ctx).Save(book).Error
}

func (r *bookRepo) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&models.Book{}).Error
}

func (r *bookRepo) GetByID(ctx context.Context, id string) (*models.Book, error) {
	var b models.Book
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *bookRepo) List(ctx context.Context) ([]models.Book, error) {
	var list []models.Book
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Новый метод для borrows
func (r *bookRepo) CreateBorrow(ctx context.Context, borrow *models.Borrow) error {
	return r.db.WithContext(ctx).Create(borrow).Error
}
