package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/account/schemes"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

// ListAccounts returns all accounts
// @Tags accounts
// @Summary List accounts
// @Description List all accounts.
// @Accept json
// @Produce json
// @Success 200 {array} accountResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [GET]
func (h Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, responses.ErrInternalError)

		return
	}

	accountsResponse := make([]schemes.AccountListResponse, 0)

	for _, acc := range accounts {
		accountsResponse = append(accountsResponse, convertAccountToAccountListResponse(acc))
	}

	responses.SendJSON(w, http.StatusOK, accountsResponse)
}
