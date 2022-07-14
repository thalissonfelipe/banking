package account

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// GetAccountBalance returns a balance by account ID.
// @Tags Accounts
// @Summary Gets account balance.
// @Description Gets account balance by account ID, if exists.
// @Accept json
// @Produce json
// @Param accountID path string true "Account ID"
// @Success 200 {object} schema.BalanceResponse
// @Failure 400 {object} rest.Error
// @Failure 500 {object} rest.InternalServerError
// @Router /accounts/{accountID}/balance [GET].
func (h Handler) GetAccountBalance(r *http.Request) rest.Response {
	accountID, err := rest.ParseUUID(chi.URLParam(r, "accountID"), "path.accountID")
	if err != nil {
		return rest.BadRequest(err, "invalid path parameters")
	}

	balance, err := h.usecase.GetAccountBalanceByID(r.Context(), vos.AccountID(accountID))
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return rest.NotFound(err, "account not found")
		}

		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToBalanceResponse(balance))
}
