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
// @Success 200 schema.ListAccountsResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [GET].
func (h Handler) ListAccounts(r *http.Request) rest.Response {
	accounts, err := h.usecase.ListAccounts(r.Context())
	if err != nil {
		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToListAccountsResponse(accounts))
}
