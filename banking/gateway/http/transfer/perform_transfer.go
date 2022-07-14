package transfer

import (
	"errors"
	"net/http"
	"strings"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

// CreateTransfer creates a new transfer between two accounts
// @Tags transfers
// @Summary Create a new transfer
// @Description Creates a new transfer between two accounts.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Param Body body transferRequest true "Body"
// @Success 201 schema.PerformTransferResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /transfers [POST].
func (h Handler) PerformTransfer(r *http.Request) rest.Response {
	var body schema.PerformTransferInput
	if err := rest.DecodeRequestBody(r, &body); err != nil {
		return rest.BadRequest(err, "invalid request body")
	}

	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	accountID, err := rest.ParseUUID(jwt.GetAccountIDFromToken(token), "token")
	if err != nil {
		return rest.BadRequest(err, "invalid token")
	}

	accountDestinationID, err := rest.ParseUUID(body.AccountDestinationID, "body.account_destination_id")
	if err != nil {
		return rest.BadRequest(err, "invalid request body")
	}

	if accountID == accountDestinationID {
		return rest.BadRequest(rest.ErrSameAccounts, rest.ErrSameAccounts.Error())
	}

	input := usecases.NewPerformTransferInput(vos.AccountID(accountID), vos.AccountID(accountDestinationID), body.Amount)

	err = h.usecase.PerformTransfer(r.Context(), input)
	if err != nil {
		if errors.Is(err, entity.ErrInsufficientFunds) {
			return rest.BadRequest(err, "insufficient funds")
		}

		if errors.Is(err, entity.ErrAccountNotFound) {
			return rest.NotFound(err, "account origin not found")
		}

		if errors.Is(err, entity.ErrAccountDestinationNotFound) {
			return rest.NotFound(err, "account destination not found")
		}

		return rest.InternalServer(err)
	}

	return rest.Created(schema.MapToPerformTransferResponse(accountID.String(), body.AccountDestinationID, body.Amount))
}
