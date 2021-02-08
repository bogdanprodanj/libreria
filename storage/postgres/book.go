package postgres

import (
	"context"
	"fmt"

	"github.com/libreria/models"
)

func (s *Storage) CreateBook(ctx context.Context, b *models.Book) error {
	_, err := s.db.WithContext(ctx).Model(b).Insert()
	return toServiceError(err)
}

func (s *Storage) UpdateBook(ctx context.Context, b *models.Book) error {
	res, err := s.db.WithContext(ctx).Model(b).WherePK().
		Set("title = ?title").
		Set("author = ?author").
		Set("publisher = ?publisher").
		Set("publish_date = ?publish_date").
		Update()
	if err != nil {
		return toServiceError(err)
	}
	if res.RowsAffected() == 0 {
		return models.ErrNotFound{Message: "book does not exist"}
	}
	return nil
}

func (s *Storage) GetBook(ctx context.Context, id int) (*models.Book, error) {
	var res models.Book
	err := s.db.WithContext(ctx).Model(&res).Where("id = ?", id).First()
	if err != nil {
		return nil, toServiceError(err)
	}
	return &res, nil
}

func (s *Storage) GetBooks(ctx context.Context, bs *models.BookSearch, limit, offset int) ([]models.Book, error) {
	var res []models.Book
	q := s.db.WithContext(ctx).Model(&res).Limit(limit).Offset(offset)
	if bs != nil {
		if bs.Title != "" {
			q = q.Where("title LIKE ?", "%"+bs.Title+"%")
		}
		if bs.Author != "" {
			q = q.Where("author LIKE ?", "%"+bs.Author+"%")
		}
		if bs.Publisher != "" {
			q = q.Where("publisher LIKE ?", "%"+bs.Publisher+"%")
		}
		if bs.Status != nil {
			q = q.Where("status = ?", *bs.Status)
		}
		if bs.PublishDateSearch != nil {
			q = q.Where(fmt.Sprintf("publish_date %s ?::date", bs.PublishDateSearch.Condition), bs.PublishDateSearch.PublishDate)
		}
	}
	err := q.Select()
	if err != nil {
		return nil, toServiceError(err)
	}
	return res, nil
}

func (s *Storage) UpdateBookStatus(ctx context.Context, id, status int) error {
	_, err := s.db.WithContext(ctx).Model((*models.Book)(nil)).
		Set("status = ?", status).Where("id = ?", id).Update()
	return toServiceError(err)
}

func (s *Storage) RateBook(ctx context.Context, id, rate int) error {
	_, err := s.db.WithContext(ctx).Model((*models.Book)(nil)).Exec(`
UPDATE ?TableName
SET rating = avg_rating
FROM (
         SELECT id,
                CASE
                    WHEN rating = 0 THEN 3.0
                    WHEN rating > 0 THEN ROUND(CAST((? + rating) / 2 AS NUMERIC), 2) END avg_rating
         FROM ?TableName
         WHERE id = ?) AS subquery
WHERE books.id = subquery.id`, rate, id)
	return toServiceError(err)
}

func (s *Storage) DeleteBook(ctx context.Context, id int) error {
	_, err := s.db.WithContext(ctx).Model((*models.Book)(nil)).Where("id = ?", id).Delete()
	return toServiceError(err)
}
