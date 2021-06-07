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
	errDestIDEqualCurrentID      = errors.New("account destination cannot be the account origin id")
)

type TransferListResponse struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
	CreatedAt            string `json:"created_at"`
}

type CreateTransferInput struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

type CreateTransferResponse struct {
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

func (t CreateTransferInput) isValid() error {
	if t.AccountDestinationID == "" {
		return errMissingAccDestIDParameter
	}
	if t.Amount == 0 {
		return errMissingAmountParameter
	}

	return nil
}

func convertTransferToTransferListResponse(transfer entities.Transfer) TransferListResponse {
	return TransferListResponse{
		AccountDestinationID: transfer.AccountDestinationID.String(),
		Amount:               transfer.Amount,
		CreatedAt:            formatTime(transfer.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}

func getTokenFromHeader(authHeader string) string {
	splitToken := strings.Split(authHeader, "Bearer ")
	return splitToken[1]
}
