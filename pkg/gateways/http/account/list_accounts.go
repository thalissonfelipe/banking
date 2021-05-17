package account

import (
	"encoding/json"
	"net/http"
)

func (h Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("Internal Error."))
	}

	var accountsResponse []AccountResponse
	for _, acc := range accounts {
		accountsResponse = append(accountsResponse, convertAccountToAccountResponse(acc))
	}

	json.NewEncoder(w).Encode(accountsResponse)
}
