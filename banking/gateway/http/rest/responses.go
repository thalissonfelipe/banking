package rest

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

type Response struct {
	Status  int
	Payload interface{}
	Error   error
}

type Error struct {
	Error   string        `json:"error" example:"invalid request body"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorDetail struct {
	Location string `json:"location" example:"body.name"`
	Message  string `json:"message" example:"missing parameter"`
}

func OK(payload interface{}) Response {
	return Response{
		Status:  http.StatusOK,
		Payload: payload,
	}
}

func Created(payload interface{}) Response {
	return Response{
		Status:  http.StatusCreated,
		Payload: payload,
	}
}

func BadRequest(err error, msg string) Response {
	return Response{
		Status:  http.StatusBadRequest,
		Payload: badRequest(err, msg),
		Error:   err,
	}
}

// TODO: add tests.
func badRequest(err error, msg string) Error {
	const bodyPrefix = "body."

	berr := Error{
		Error:   msg,
		Details: []ErrorDetail{},
	}

	if errors.Is(err, vos.ErrInvalidCPF) {
		berr.Details = append(berr.Details, ErrorDetail{
			Message:  vos.ErrInvalidCPF.Error(),
			Location: bodyPrefix + "cpf",
		})
	}

	if errors.Is(err, vos.ErrInvalidSecret) {
		berr.Details = append(berr.Details, ErrorDetail{
			Message:  vos.ErrInvalidSecret.Error(),
			Location: bodyPrefix + "secret",
		})
	}

	var verr ValidationError
	if errors.As(err, &verr) {
		berr.Details = append(berr.Details, ErrorDetail{
			Message:  verr.Err.Error(),
			Location: verr.Location,
		})
	}

	var verrs ValidationErrors
	if errors.As(err, &verrs) {
		for _, err := range verrs {
			var verr ValidationError
			if errors.As(err, &verr) {
				berr.Details = append(berr.Details, ErrorDetail{
					Message:  verr.Err.Error(),
					Location: verr.Location,
				})
			}
		}
	}

	return berr
}

func InvalidCredentials(err error) Response {
	return Response{
		Status: http.StatusBadRequest,
		Payload: Error{
			Error: "invalid credentials",
		},
		Error: err,
	}
}

func Unauthorized(err error) Response {
	return Response{
		Status: http.StatusUnauthorized,
		Payload: Error{
			Error: "unauthorized",
		},
		Error: err,
	}
}

func NotFound(err error, msg string) Response {
	return Response{
		Status: http.StatusNotFound,
		Payload: Error{
			Error: msg,
		},
		Error: err,
	}
}

func Conflict(err error, msg string) Response {
	return Response{
		Status: http.StatusConflict,
		Payload: Error{
			Error: msg,
		},
		Error: err,
	}
}

func InternalServer(err error) Response {
	return Response{
		Status: http.StatusInternalServerError,
		Payload: Error{
			Error: "internal server error",
		},
		Error: err,
	}
}
