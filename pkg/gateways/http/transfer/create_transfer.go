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

func (h Handler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	var body transferRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, errInvalidJSON.Error())
		return
	}

	if err := body.isValid(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	token := getTokenFromHeader(r.Header.Get("Authorization"))
	accountID := vos.ConvertStringToID(auth.GetIDFromToken(token))
	accountDestinationID := vos.ConvertStringToID(body.AccountDestinationID)

	if accountID == accountDestinationID {
		responses.SendError(w, http.StatusBadRequest, errDestIDEqualCurrentID.Error())
		return
	}

	input := transfer.NewTransferInput(accountID, accountDestinationID, body.Amount)
	err = h.usecase.CreateTransfer(r.Context(), input)
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			responses.SendError(w, http.StatusNotFound, errAccountOriginDoesNotExist.Error())
			return
		}
		if errors.Is(err, entities.ErrAccountDestinationDoesNotExist) {
			responses.SendError(w, http.StatusNotFound, err.Error())
			return
		}
		if errors.Is(err, entities.ErrInsufficientFunds) {
			responses.SendError(w, http.StatusBadRequest, err.Error())
			return
		}

		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := transferCreatedResponse{
		AccountOriginID:      accountID.String(),
		AccountDestinationID: body.AccountDestinationID,
		Amount:               body.Amount,
	}
	responses.SendJSON(w, http.StatusOK, response)
}
