package transfer

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/rest"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/transfer/schemes"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

// CreateTransfer creates a new transfer between two accounts
// @Tags transfers
// @Summary Create a new transfer
// @Description Creates a new transfer between two accounts.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Param Body body transferRequest true "Body"
// @Success 201 {array} transferCreatedResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /transfers [POST].
func (h Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var body schemes.CreateTransferInput

	if err := rest.DecodeRequestBody(r, &body); err != nil {
		rest.HandleBadRequestError(w, err)

		return
	}

	token := getTokenFromHeader(r.Header.Get("Authorization"))
	accountID := vos.ConvertStringToAccountID(auth.GetIDFromToken(token))
	accountDestinationID := vos.ConvertStringToAccountID(body.AccountDestinationID)

	if accountID == accountDestinationID {
		rest.SendError(w, http.StatusBadRequest, rest.ErrDestinationIDEqToOriginID)

		return
	}

	input := transfer.NewTransferInput(accountID, accountDestinationID, body.Amount)

	err := h.usecase.CreateTransfer(r.Context(), input)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			rest.SendError(w, http.StatusNotFound, rest.ErrAccountOriginNotFound)

			return
		}

		rest.HandleError(w, err)

		return
	}

	response := schemes.CreateTransferResponse{
		AccountOriginID:      accountID.String(),
		AccountDestinationID: body.AccountDestinationID,
		Amount:               body.Amount,
	}
	rest.SendJSON(w, http.StatusCreated, response)
}
