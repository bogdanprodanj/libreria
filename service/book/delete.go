package book

import "context"

func (s *Service) DeleteBook(ctx context.Context, id int) error {
	return s.storage.DeleteBook(ctx, id)
}
