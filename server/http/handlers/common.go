package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/libreria/models"
	hm "github.com/libreria/server/http/models"
	log "github.com/sirupsen/logrus"
)

// unmarshalRequest unmarshalls http request body to provided structure.
func unmarshalRequestBody(r *http.Request, body interface{}) error {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return models.ErrInternal{Message: err.Error()}
	}
	defer r.Body.Close()
	if err := json.Unmarshal(reqBody, body); err != nil {
		return models.ErrBadRequest{Message: err.Error()}
	}
	return nil
}

const defaultLimit = 50

func getPagination(r *http.Request) (limit, offset int, err error) {
	limit = defaultLimit
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	if limitStr != "" {
		l, lErr := strconv.Atoi(limitStr)
		if lErr != nil {
			err = models.ErrBadRequest{Message: "invalid limit format"}
		}
		limit = l
	}
	if offsetStr != "" {
		o, oErr := strconv.Atoi(offsetStr)
		if oErr != nil {
			err = models.ErrBadRequest{Message: "invalid offset format"}
		}
		offset = o
	}
	return
}

const timeFormat = "2006-01-02"

func getSearch(r *http.Request) (*models.BookSearch, error) {
	status := r.URL.Query().Get("status")
	var intStatus *int
	if status != "" {
		if status == strings.ToLower(string(hm.StatusCheckedIn)) {
			intStatus = toIntPtr(0)
		} else if status == strings.ToLower(string(hm.StatusCheckedOut)) {
			intStatus = toIntPtr(1)
		} else {
			return nil, models.ErrBadRequest{Message: "invalid status format, use 'checkedIn' or 'checkedOut'"}
		}
	}
	var pds *models.PublishDateSearch
	pd := r.URL.Query().Get("publish_date")
	if pd != "" {
		filter := strings.Split(pd, " ")
		if len(filter) != 2 {
			return nil, models.ErrBadRequest{Message: "invalid publish_date filter format, use 'lte yyyy-dd-mm"}
		}
		ft, ok := models.FilterMap[filter[0]]
		if !ok {
			return nil, models.ErrBadRequest{Message: "invalid publish_date filter format, use 'lte yyyy-dd-mm"}
		}
		t, err := time.Parse(timeFormat, filter[1])
		if err != nil {
			return nil, models.ErrBadRequest{Message: "invalid publish_date filter format, use 'lte yyyy-dd-mm"}
		}
		pds = &models.PublishDateSearch{
			PublishDate: t.Format(timeFormat),
			Condition:   ft,
		}
	}
	return &models.BookSearch{
		Title:             r.URL.Query().Get("title"),
		Author:            r.URL.Query().Get("author"),
		Publisher:         r.URL.Query().Get("publisher"),
		Status:            intStatus,
		PublishDateSearch: pds,
	}, nil

}

func toIntPtr(b int) (ptr *int) {
	ptr = &b
	return
}

// sendEmptyResponse sends only response code
func sendEmptyResponse(w http.ResponseWriter, statusCode int) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(statusCode)
}

// sendResponseWithBody sends response code and body in JSON format
func sendResponseWithBody(w http.ResponseWriter, statusCode int, respBody interface{}) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	b, err := json.Marshal(respBody)
	if err != nil {
		statusCode = http.StatusInternalServerError
	}
	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

func sendInvalidIDError(w http.ResponseWriter) {
	sendHTTPError(w, &models.ErrBadRequest{Message: "invalid book id"})
}

// sendHTTPError sends error response with appropriate status code.
func sendHTTPError(w http.ResponseWriter, err error) {
	var (
		code    int
		message string
		errs    []models.FieldError
	)
	switch v := err.(type) {
	case validation.Error:
		code = http.StatusBadRequest
		message = v.Message()
	case validation.Errors:
		code = http.StatusBadRequest
		message = v.Error()
	case models.ErrInternal:
		log.WithError(err).Error("internal error")
		code = http.StatusInternalServerError
		message = "oops, something went wrong"
	case models.ErrNotFound:
		code = http.StatusNotFound
		message = v.Message
	case models.ErrBadRequest:
		code = http.StatusBadRequest
		message = v.Message
		errs = v.Errors
	default:
		log.WithError(err).Error("unknown error")
		code = http.StatusServiceUnavailable
		message = "service unavailable"
	}
	sendResponseWithBody(w, code, struct {
		Code    int                 `json:"code"`
		Message string              `json:"message,omitempty"`
		Errors  []models.FieldError `json:"errors,omitempty"`
	}{
		Code:    code,
		Message: message,
		Errors:  errs,
	})
}
