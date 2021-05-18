package account

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := mux.Vars(r)["id"]
	balance, err := h.usecase.GetAccountBalanceByID(r.Context(), accountID)
	if err != nil {
		switch err {
		case entities.ErrAccountDoesNotExist:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Account not found."))
			return
		default:
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Error."))
			return
		}
	}

	response := BalanceResponse{Balance: balance}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&response)
}
