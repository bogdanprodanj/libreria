package models

type apiError struct {
	Message string `json:"message"`
	Errors  []FieldError
}

type FieldError struct {
	Field string
}

type ErrNotFound apiError

func (e ErrNotFound) Error() string {
	return e.Message
}

type ErrInternal apiError

func (e ErrInternal) Error() string {
	return e.Message
}

type ErrBadRequest apiError

func (e ErrBadRequest) Error() string {
	return e.Message
}
