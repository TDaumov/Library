package service

import (
	"context"
	"library_vebservice/models"
	"library_vebservice/repository"
	//"github.com/yourname/library-api/internal/models"
	//"github.com/yourname/library-api/internal/repository"
)

type BookService interface {
	List(ctx context.Context) ([]models.Book, error)
	GetByID(ctx context.Context, id string) (*models.Book, error)
	Create(ctx context.Context, book *models.Book) error
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, id string) error
	CreateBorrow(ctx context.Context, borrow *models.Borrow) error // ← добавляем
}

type bookService struct {
	repo repository.BookRepo
}

func NewBookService(r repository.BookRepo) BookService {
	return &bookService{repo: r}
}

func (s *bookService) CreateBorrow(ctx context.Context, borrow *models.Borrow) error {
	return s.repo.CreateBorrow(ctx, borrow)
}

func (s *bookService) Create(ctx context.Context, b *models.Book) error { return s.repo.Create(ctx, b) }
func (s *bookService) Update(ctx context.Context, b *models.Book) error { return s.repo.Update(ctx, b) }
func (s *bookService) Delete(ctx context.Context, id string) error      { return s.repo.Delete(ctx, id) }
func (s *bookService) GetByID(ctx context.Context, id string) (*models.Book, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *bookService) List(ctx context.Context) ([]models.Book, error) { return s.repo.List(ctx) }
