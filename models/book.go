package models

import "time"

type Book struct {
	ID          int       `json:"id" pg:",pk"`
	Title       string    `json:"title" pg:"title"`
	Author      string    `json:"author" pg:"author"`
	Publisher   string    `json:"publisher" pg:"publisher"`
	PublishDate time.Time `json:"publish_date" pg:"publish_date"`
	Rating      float64   `json:"rating" pg:"rating"`
	Status      int       `json:"status" pg:"status"`
	DeletedAt   *time.Time `pg:",soft_delete" json:"-" `
	CreatedAt   *time.Time `pg:"default:now()" json:"-" `
	UpdatedAt   *time.Time `pg:"default:now()" json:"-" `
}

type BookSearch struct {
	Title             string
	Author            string
	Publisher         string
	Status            *int
	PublishDateSearch *PublishDateSearch
}

type PublishDateSearch struct {
	PublishDate string
	Condition   string
}

var FilterMap = map[string]string{
	"eq":  "=",
	"neq": "<>",
	"lt":  "<",
	"lte": "<=",
	"gt":  ">",
	"gte": ">=",
}
