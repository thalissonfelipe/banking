package account

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

func (h Handler) GetAccountBalance(w http.ResponseWriter, r *http.Request) {
	accountID := vos.ConvertStringToID(mux.Vars(r)["id"])
	balance, err := h.usecase.GetAccountBalanceByID(r.Context(), accountID)
	if err != nil {
		switch err {
		case entities.ErrAccountDoesNotExist:
			responses.SendError(w, http.StatusNotFound, err.Error())
		default:
			responses.SendError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response := balanceResponse{Balance: balance}
	responses.SendJSON(w, http.StatusOK, response)
}
