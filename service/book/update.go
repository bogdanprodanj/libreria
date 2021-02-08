package book

import (
	"context"

	"github.com/libreria/models"
)

func (s *Service) UpdateBook(ctx context.Context, b *models.Book) error {
	return s.storage.UpdateBook(ctx, b)
}

func (s *Service) UpdateBookStatus(ctx context.Context, id, status int) error {
	return s.storage.UpdateBookStatus(ctx, id, status)
}

func (s *Service) RateBook(ctx context.Context, id, rate int) error {
	return s.storage.RateBook(ctx, id, rate)
}
