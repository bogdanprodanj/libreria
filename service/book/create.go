package book

import (
	"context"

	"github.com/libreria/models"
)

func (s *Service) AddBook(ctx context.Context, b *models.Book) error {
	return s.storage.CreateBook(ctx, b)
}
