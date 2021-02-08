package book

import (
	"context"

	"github.com/libreria/models"
)

func (s *Service) GetBook(ctx context.Context, id int) (*models.Book, error) {
	return s.storage.GetBook(ctx, id)
}

func (s *Service) GetBooks(ctx context.Context, bs *models.BookSearch, limit, offset int) ([]models.Book, error) {
	return s.storage.GetBooks(ctx, bs, limit, offset)
}
