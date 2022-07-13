package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/services/auth"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func HandleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, entities.ErrInsufficientFunds):
		SendError(w, http.StatusBadRequest, ErrInsufficientFunds)
	case errors.Is(err, auth.ErrInvalidCredentials):
		SendError(w, http.StatusBadRequest, ErrInvalidCredentials)
	case errors.Is(err, entities.ErrAccountDoesNotExist):
		SendError(w, http.StatusNotFound, ErrAccountNotFound)
	case errors.Is(err, entities.ErrAccountDestinationDoesNotExist):
		SendError(w, http.StatusNotFound, ErrAccountDestinationNotFound)
	case errors.Is(err, entities.ErrAccountAlreadyExists):
		SendError(w, http.StatusConflict, ErrAccountAlreadyExists)
	default:
		SendError(w, http.StatusInternalServerError, ErrInternalError)
	}
}

func HandleBadRequestError(w http.ResponseWriter, err error) {
	const status = http.StatusBadRequest

	switch {
	case errors.Is(err, ErrInvalidJSON):
		SendError(w, status, ErrInvalidJSON)
	case errors.Is(err, ErrMissingNameParameter):
		SendError(w, status, ErrMissingNameParameter)
	case errors.Is(err, ErrMissingCPFParameter):
		SendError(w, status, ErrMissingCPFParameter)
	case errors.Is(err, ErrMissingSecretParameter):
		SendError(w, status, ErrMissingSecretParameter)
	case errors.Is(err, ErrMissingAccDestinationIDParameter):
		SendError(w, status, ErrMissingAccDestinationIDParameter)
	case errors.Is(err, ErrMissingAmountParameter):
		SendError(w, status, ErrMissingAmountParameter)
	}
}

func SendError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
}

func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	_ = json.NewEncoder(w).Encode(payload)
}
