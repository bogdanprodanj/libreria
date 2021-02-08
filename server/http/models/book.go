package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Book struct {
	Title       string    `json:"name"`
	Author      string    `json:"author"`
	Publisher   string    `json:"publisher"`
	PublishDate time.Time `json:"publish_date"`
}

func (b Book) Validate() error {
	return validation.ValidateStruct(&b,
		validation.Field(&b.Title, validation.Required, validation.Length(1, 200)),
		validation.Field(&b.Author, validation.Required, validation.Length(1, 200)),
		validation.Field(&b.Publisher, validation.Required, validation.Length(1, 200)),
		validation.Field(&b.PublishDate, validation.Required, validation.Max(time.Now())),
	)
}

type Status string

const (
	StatusCheckedIn  Status = "CheckedIn"
	StatusCheckedOut Status = "CheckedOut"
)

type GetBookResponse struct {
	Book
	ID     int     `json:"id"`
	Status Status  `json:"status"`
	Rating float64 `json:"rating,omitempty"`
}

type RateRequest struct {
	Rating int `json:"rating"`
}

func (r RateRequest) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Rating, validation.Required, validation.Min(1), validation.Max(3)),
	)
}
