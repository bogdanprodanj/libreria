// +build unit

package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/libreria/models"
	"github.com/libreria/server/http/handlers/mock"
	hm "github.com/libreria/server/http/models"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBook_UpdateBook(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	a := assert.New(t)    // assertion object for comparing values
	req := require.New(t) // same as assertion, but stops test execution if condition is false

	ctrl := gomock.NewController(t) // gomock controller
	defer ctrl.Finish()
	srvMock := mock.NewMockBookKeeper(ctrl) // mocked service
	oh := New(srvMock)                      // book handler with mocked service

	router := mux.NewRouter()                                               // router
	router.HandleFunc("/books/{id}", oh.UpdateBook).Methods(http.MethodPut) // create book route
	srv := httptest.NewServer(router)                                       // test server
	defer srv.Close()
	input := &hm.Book{
		Title:       "my_title",
		Author:      "my_author",
		Publisher:   "my_publisher",
		PublishDate: time.Now().Add(-time.Hour).UTC(),
	}
	book := &models.Book{
		ID:          1,
		Title:       input.Title,
		Author:      input.Author,
		Publisher:   input.Publisher,
		PublishDate: input.PublishDate,
	}
	t.Run("happy_path", func(t *testing.T) {
		srvMock.EXPECT().UpdateBook(gomock.Any(), book).Return(nil)

		reqBody, err := json.Marshal(input)
		req.NoError(err)

		request, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/books/1", srv.URL), bytes.NewBuffer(reqBody))
		req.NoError(err)
		res, err := srv.Client().Do(request)
		req.NoError(err)
		defer res.Body.Close()
		a.Equal(http.StatusNoContent, res.StatusCode)
	})
}

func TestBook_AddBook(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	a := assert.New(t)    // assertion object for comparing values
	req := require.New(t) // same as assertion, but stops test execution if condition is false

	ctrl := gomock.NewController(t) // gomock controller
	defer ctrl.Finish()
	srvMock := mock.NewMockBookKeeper(ctrl) // mocked service
	oh := New(srvMock)                      // book handler with mocked service

	router := mux.NewRouter()                                        // router
	router.HandleFunc("/books", oh.AddBook).Methods(http.MethodPost) // create book route
	srv := httptest.NewServer(router)                                // test server
	defer srv.Close()
	input := &hm.Book{
		Title:       "my_title",
		Author:      "my_author",
		Publisher:   "my_publisher",
		PublishDate: time.Now().Add(-time.Hour).UTC(),
	}
	book := &models.Book{
		Title:       input.Title,
		Author:      input.Author,
		Publisher:   input.Publisher,
		PublishDate: input.PublishDate,
	}
	t.Run("happy_path", func(t *testing.T) {
		srvMock.EXPECT().AddBook(gomock.Any(), book).Return(nil)

		reqBody, err := json.Marshal(input)
		req.NoError(err)

		res, err := http.Post(fmt.Sprintf("%s/books", srv.URL), "application/json", bytes.NewBuffer(reqBody))
		req.NoError(err)
		defer res.Body.Close()
		a.Equal(http.StatusCreated, res.StatusCode)

		body, err := ioutil.ReadAll(res.Body)
		req.NoError(err)

		var bookFromHandler hm.GetBookResponse
		err = json.Unmarshal(body, &bookFromHandler)
		req.NoError(err)

		a.Equal(book.Title, bookFromHandler.Title)
		a.Equal(hm.StatusCheckedIn, bookFromHandler.Status)
	})
	t.Run("validation", func(t *testing.T) {
		input := &hm.Book{PublishDate: time.Now().Add(-time.Hour).UTC()}
		reqBody, err := json.Marshal(input)
		req.NoError(err)
		resp, err := http.Post(fmt.Sprintf("%s/books", srv.URL), "application/json", bytes.NewBuffer(reqBody))
		req.NoError(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode) // Check if the Code is 400
	})
	t.Run("bad_request", func(t *testing.T) {
		invalid := "adc" // invalid request body
		resp, err := http.Post(fmt.Sprintf("%s/books", srv.URL), "application/json", bytes.NewBuffer([]byte(invalid)))
		req.NoError(err)
		a.Equal(http.StatusBadRequest, resp.StatusCode) // Check if the Code is 400
	})
	t.Run("service_unavailable", func(t *testing.T) {
		errToReturn := errors.New("internal")
		srvMock.EXPECT().AddBook(gomock.Any(), book).Return(errToReturn) // Expect database to return errToReturn
		reqBody, err := json.Marshal(input)
		req.NoError(err)
		resp, err := http.Post(fmt.Sprintf("%s/books", srv.URL), "application/json", bytes.NewBuffer(reqBody))
		req.NoError(err)
		a.Equal(http.StatusServiceUnavailable, resp.StatusCode) // Check if the Code is 503
	})
}
