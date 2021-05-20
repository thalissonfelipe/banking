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
)

func getTokenFromHeader(authHeader string) string {
	splitToken := strings.Split(authHeader, "Bearer ")
	return splitToken[1]
}

type transferResponse struct {
	AccountDestinationID string
	Amount               int
	CreatedAt            string
}

type transferRequest struct {
	AccountDestinationID string
	Amount               int
}

type transferCreatedResponse struct {
	AccountOriginID      string
	AccountDestinationID string
	Amount               int
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
		AccountDestinationID: transfer.AccountDestinationID,
		Amount:               transfer.Amount,
		CreatedAt:            formatTime(transfer.CreatedAt),
	}
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000000")
}
