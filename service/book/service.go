package book

import (
	"context"

	"github.com/libreria/models"
)

type StorageManager interface {
	GetBook(ctx context.Context, id int) (*models.Book, error)
	GetBooks(ctx context.Context, bs *models.BookSearch, limit, offset int) ([]models.Book, error)
	UpdateBook(ctx context.Context, b *models.Book) error
	CreateBook(ctx context.Context, b *models.Book) error
	UpdateBookStatus(ctx context.Context, id, status int) error
	RateBook(ctx context.Context, id, rate int) error
	DeleteBook(ctx context.Context, id int) error
}

type Service struct {
	storage StorageManager
}

func New(storage StorageManager) *Service {
	return &Service{storage: storage}
}
