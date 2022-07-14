package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// ListAccounts returns all accounts.
// @Tags Accounts
// @Summary Lists all accounts.
// @Description Lists all accounts.
// @Accept json
// @Produce json
// @Success 200 {object} schema.ListAccountsResponse
// @Failure 500 {object} rest.InternalServerError
// @Router /accounts [GET].
func (h Handler) ListAccounts(r *http.Request) rest.Response {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToListAccountsResponse(accounts))
}
