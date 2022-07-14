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

// GetAccountBalance returns a balance by accountID
// @Tags accounts
// @Summary Get account balance
// @Description Get account balance by accountID, if exists.
// @Accept json
// @Produce json
// @Success 200 {object} schema.BalanceResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
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
