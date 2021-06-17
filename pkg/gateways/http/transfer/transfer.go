package transfer

import (
	"strings"
	"time"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/transfer/schemes"
)

func convertTransferToTransferListResponse(transfer entities.Transfer) schemes.TransferListResponse {
	return schemes.TransferListResponse{
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
