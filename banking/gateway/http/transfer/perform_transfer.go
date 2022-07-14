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

// PerformTransfer performs a transfer between two accounts.
// @Tags Transfers
// @Summary Performs a transfer between two accounts.
// @Description Performs a transfer between two accounts. User must be authenticated.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Param Body body schema.PerformTransferInput true "Body"
// @Success 201 {object} schema.PerformTransferResponse
// @Failure 400 {object} rest.Error
// @Failure 401 {object} rest.UnauthorizedError
// @Failure 404 {object} rest.NotFoundError
// @Failure 500 {object} rest.InternalServerError
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
