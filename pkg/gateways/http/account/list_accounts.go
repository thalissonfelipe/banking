package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

func (h Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accountsResponse := make([]accountResponse, 0)
	for _, acc := range accounts {
		accountsResponse = append(accountsResponse, convertAccountToAccountResponse(acc))
	}

	responses.SendJSON(w, http.StatusOK, accountsResponse)
}
