package responses

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

var (
	ErrInvalidJSON                = errors.New("invalid json")
	ErrAccountNotFound            = errors.New("account does not exist")
	ErrAccountOriginNotFound      = errors.New("account origin does not exist")
	ErrAccountDestinationNotFound = errors.New("account destination does not exist")
	ErrInternalError              = errors.New("internal server error")
	ErrInsufficientFunds          = errors.New("insufficient funds")
	ErrInvalidCredentials         = errors.New("cpf or secret are invalid")
	ErrAccountAlreadyExists       = errors.New("account already exists")
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
	case errors.Is(err, entities.ErrInternalError):
		SendError(w, http.StatusInternalServerError, ErrInternalError)
	}
}

func SendError(w http.ResponseWriter, statusCode int, err error) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Message: err.Error()})
}

func SendJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
