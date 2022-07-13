package transfer

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/services/auth"
)

// CreateTransfer creates a new transfer between two accounts
// @Tags transfers
// @Summary Create a new transfer
// @Description Creates a new transfer between two accounts.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Param Body body transferRequest true "Body"
// @Success 201 {array} schema.PerformTransferResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /transfers [POST].
func (h Handler) PerformTransfer(w http.ResponseWriter, r *http.Request) {
	var body schema.PerformTransferInput

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

	input := usecases.NewPerformTransferInput(accountID, accountDestinationID, body.Amount)

	err := h.usecase.PerformTransfer(r.Context(), input)
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			rest.SendError(w, http.StatusNotFound, rest.ErrAccountOriginNotFound)

			return
		}

		rest.HandleError(w, err)

		return
	}

	response := schema.MapToPerformTransferResponse(accountID.String(), body.AccountDestinationID, body.Amount)
	rest.SendJSON(w, http.StatusCreated, response)
}
