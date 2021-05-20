package transfer

import (
	"strings"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
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
