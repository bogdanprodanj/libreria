// +build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"

	hm "github.com/libreria/server/http/models"
)

func (s *LibreriaTestSuite) TestGetBooks() {
	s.Run("all_books", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var books []hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Assert().Len(books, 3)
	})
	s.Run("limit_1", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books?limit=1", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var books []hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Assert().Len(books, 1)
	})
	s.Run("publish_date", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books?publish_date=gt%202020-01-01", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var books []hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Assert().Len(books, 1)
	})
	s.Run("title_contains", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books?title=Faster", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var books []hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Len(books, 1)
		s.Assert().Equal(testBooks[1].Title, books[0].Title)
	})
	s.Run("no_books", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books?title=Hobbit", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var books []hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&books)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Len(books, 0)
	})
}

func (s *LibreriaTestSuite) TestGetBook() {
	s.Run("found", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/1", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var book hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&book)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(testBooks[0].Title, book.Title)
	})
	s.Run("not_found", func() {
		req, _ := http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/100", nil)
		resp, err := s.c.Do(req)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *LibreriaTestSuite) TestCheckInBook() {
	s.Run("found", func() {
		req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8080/api/v1/books/1/out", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
		req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/1", nil)
		resp, err = s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var book hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&book)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(hm.StatusCheckedOut, book.Status)
	})
}

func (s *LibreriaTestSuite) TestDeleteBook() {
	s.Run("happy_path", func() {
		req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8080/api/v1/books/1", nil)
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
		req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/1", nil)
		resp, err = s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNotFound, resp.StatusCode)
	})
}

func (s *LibreriaTestSuite) TestRateBook() {
	s.Run("happy_path", func() {
		rate := &hm.RateRequest{Rating: 3}
		b, _ := json.Marshal(rate)
		req, _ := http.NewRequest(http.MethodPatch, "http://localhost:8080/api/v1/books/1/rate", bytes.NewReader(b))
		resp, err := s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
		req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/1", nil)
		resp, err = s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var book hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&book)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(3.0, book.Rating)

		rate.Rating = 2
		b, _ = json.Marshal(rate)
		req, _ = http.NewRequest(http.MethodPatch, "http://localhost:8080/api/v1/books/1/rate", bytes.NewReader(b))
		resp, err = s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusNoContent, resp.StatusCode)
		req, _ = http.NewRequest(http.MethodGet, "http://localhost:8080/api/v1/books/1", nil)
		resp, err = s.c.Do(req)
		s.Require().NoError(err)
		s.Require().Equal(http.StatusOK, resp.StatusCode)
		var book2 hm.GetBookResponse
		err = json.NewDecoder(resp.Body).Decode(&book2)
		resp.Body.Close()
		s.Require().NoError(err)
		s.Require().Equal(2.5, book2.Rating)
	})
}
