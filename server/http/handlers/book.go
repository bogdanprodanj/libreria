package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/libreria/models"
	hm "github.com/libreria/server/http/models"
)

type BookKeeper interface {
	AddBook(ctx context.Context, b *models.Book) error
	GetBook(ctx context.Context, id int) (*models.Book, error)
	GetBooks(ctx context.Context, bs *models.BookSearch, limit, offset int) ([]models.Book, error)
	UpdateBook(ctx context.Context, b *models.Book) error
	UpdateBookStatus(ctx context.Context, id, status int) error
	RateBook(ctx context.Context, id, rate int) error
	DeleteBook(ctx context.Context, id int) error
}

type Book struct {
	bk BookKeeper
}

func New(bk BookKeeper) *Book {
	return &Book{bk: bk}
}

func (h *Book) AddBook(w http.ResponseWriter, r *http.Request) {
	var req hm.Book
	err := unmarshalRequestBody(r, &req)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	err = req.Validate()
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	book := &models.Book{
		Title:       req.Title,
		Author:      req.Author,
		Publisher:   req.Publisher,
		PublishDate: req.PublishDate.UTC(),
	}
	err = h.bk.AddBook(r.Context(), book)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	resp := toBookResponse(book)
	sendResponseWithBody(w, http.StatusCreated, &resp)
}

func (h *Book) GetBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendInvalidIDError(w)
		return
	}
	book, err := h.bk.GetBook(r.Context(), id)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	resp := toBookResponse(book)
	sendResponseWithBody(w, http.StatusOK, &resp)
}

func (h *Book) ListBooks(w http.ResponseWriter, r *http.Request) {
	limit, offset, err := getPagination(r)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	sc, err := getSearch(r)
	if err != nil {
		sendHTTPError(w, err)
		return
	}

	books, err := h.bk.GetBooks(r.Context(), sc, limit, offset)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	resp := toBooksResponse(books)
	sendResponseWithBody(w, http.StatusOK, &resp)
}

func (h *Book) UpdateBook(w http.ResponseWriter, r *http.Request) {
	var req hm.Book
	err := unmarshalRequestBody(r, &req)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	err = req.Validate()
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendInvalidIDError(w)
		return
	}
	book := &models.Book{
		ID:          id,
		Title:       req.Title,
		Author:      req.Author,
		Publisher:   req.Publisher,
		PublishDate: req.PublishDate,
	}
	err = h.bk.UpdateBook(r.Context(), book)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	sendEmptyResponse(w, http.StatusNoContent)
}

func (h *Book) CheckinBook(w http.ResponseWriter, r *http.Request) {
	h.updateBookStatus(w, r, 0)
}

func (h *Book) CheckoutBook(w http.ResponseWriter, r *http.Request) {
	h.updateBookStatus(w, r, 1)
}

func (h *Book) updateBookStatus(w http.ResponseWriter, r *http.Request, status int) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendInvalidIDError(w)
		return
	}
	err = h.bk.UpdateBookStatus(r.Context(), id, status)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	sendEmptyResponse(w, http.StatusNoContent)
}

func (h *Book) RateBook(w http.ResponseWriter, r *http.Request) {
	var req hm.RateRequest
	err := unmarshalRequestBody(r, &req)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	err = req.Validate()
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendInvalidIDError(w)
		return
	}
	err = h.bk.RateBook(r.Context(), id, req.Rating)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	sendEmptyResponse(w, http.StatusNoContent)
}

func (h *Book) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		sendInvalidIDError(w)
		return
	}
	err = h.bk.DeleteBook(r.Context(), id)
	if err != nil {
		sendHTTPError(w, err)
		return
	}
	sendEmptyResponse(w, http.StatusNoContent)
}

func toBookResponse(b *models.Book) hm.GetBookResponse {
	resp := hm.GetBookResponse{
		Book: hm.Book{
			Title:       b.Title,
			Author:      b.Author,
			Publisher:   b.Publisher,
			PublishDate: b.PublishDate,
		},
		ID:     b.ID,
		Rating: b.Rating,
	}
	if b.Status == 0 {
		resp.Status = hm.StatusCheckedIn
	} else {
		resp.Status = hm.StatusCheckedOut
	}
	return resp
}

func toBooksResponse(bs []models.Book) []hm.GetBookResponse {
	resp := make([]hm.GetBookResponse, len(bs))
	for i := range bs {
		resp[i] = toBookResponse(&bs[i])
	}
	return resp
}
