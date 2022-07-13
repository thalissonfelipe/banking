package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// ListAccounts returns all accounts
// @Tags accounts
// @Summary List accounts
// @Description List all accounts.
// @Accept json
// @Produce json
// @Success 200 {array} schema.ListAccountsResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [GET].
func (h Handler) ListAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		rest.SendError(w, http.StatusInternalServerError, rest.ErrInternalError)

		return
	}

	rest.SendJSON(w, http.StatusOK, schema.MapToListAccountsResponse(accounts))
}
