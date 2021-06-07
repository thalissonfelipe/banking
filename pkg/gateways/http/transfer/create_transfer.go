package transfer

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
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
// @Router /transfers [POST]
func (h Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var body transferRequest

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, responses.ErrInvalidJSON)

		return
	}

	if err := body.isValid(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)

		return
	}

	token := getTokenFromHeader(r.Header.Get("Authorization"))
	accountID := vos.ConvertStringToID(auth.GetIDFromToken(token))
	accountDestinationID := vos.ConvertStringToID(body.AccountDestinationID)

	if accountID == accountDestinationID {
		responses.SendError(w, http.StatusBadRequest, errDestIDEqualCurrentID)

		return
	}

	input := transfer.NewTransferInput(accountID, accountDestinationID, body.Amount)

	err = h.usecase.CreateTransfer(r.Context(), input)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			responses.SendError(w, http.StatusNotFound, responses.ErrAccountOriginNotFound)

			return
		}

		responses.HandleError(w, err)

		return
	}

	response := transferCreatedResponse{
		AccountOriginID:      accountID.String(),
		AccountDestinationID: body.AccountDestinationID,
		Amount:               body.Amount,
	}
	responses.SendJSON(w, http.StatusCreated, response)
}
