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

	accountsResponse := make([]AccountResponse, 0)
	for _, acc := range accounts {
		accountsResponse = append(accountsResponse, convertAccountToAccountResponse(acc))
	}

	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(accountsResponse)
}
