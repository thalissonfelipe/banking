package transfer

import (
	"errors"
	"strings"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
)

var (
	errMissingAccDestIDParameter = errors.New("missing account destination id parameter")
	errMissingAmountParameter    = errors.New("missing amount parameter")
	errInvalidJSON               = errors.New("invalid json")
	errAccountOriginDoesNotExist = errors.New("account origin does not exist")
	errDestIDEqualCurrentID      = errors.New("account destination cannot be the account origin id")
)

func getTokenFromHeader(authHeader string) string {
	splitToken := strings.Split(authHeader, "Bearer ")
	return splitToken[1]
}

type transferResponse struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
	CreatedAt            string `json:"created_at"`
}

type transferRequest struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

type transferCreatedResponse struct {
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

func (t transferRequest) isValid() error {
	if t.AccountDestinationID == "" {
		return errMissingAccDestIDParameter
	}
	if t.Amount == 0 {
		return errMissingAmountParameter
	}

	return nil
}

func convertTransferToTransferResponse(transfer entities.Transfer) transferResponse {
	return transferResponse{
		AccountDestinationID: transfer.AccountDestinationID.String(),
		Amount:               transfer.Amount,
		CreatedAt:            formatTime(transfer.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
