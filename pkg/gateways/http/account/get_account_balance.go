package account

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := mux.Vars(r)["id"]
	balance, err := h.usecase.GetAccountBalanceByID(r.Context(), accountID)
	if err != nil {
		switch err {
		case entities.ErrAccountDoesNotExist:
			responses.SendError(w, http.StatusNotFound, "Account not found.")
		default:
			responses.SendError(w, http.StatusInternalServerError, "Internal Error.")
		}
		return
	}

	response := BalanceResponse{Balance: balance}
	responses.SendJSON(w, http.StatusOK, response)
}
