package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/banking/gateways/http/account/schemes"
	"github.com/thalissonfelipe/banking/banking/gateways/http/rest"
)

// ListAccounts returns all accounts
// @Tags accounts
// @Summary List accounts
// @Description List all accounts.
// @Accept json
// @Produce json
// @Success 200 {array} accountResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [GET].
func (h Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		rest.SendError(w, http.StatusInternalServerError, rest.ErrInternalError)

		return
	}

	accountsResponse := make([]schemes.AccountListResponse, 0)

	for _, acc := range accounts {
		accountsResponse = append(accountsResponse, convertAccountToAccountListResponse(acc))
	}

	rest.SendJSON(w, http.StatusOK, accountsResponse)
}
