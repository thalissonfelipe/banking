package schema

import (
	"time"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

var (
	ErrMissingAccountDestIDParameter = rest.ValidationError{
		Location: "body.account_destination_id",
		Err:      rest.ErrMissingParameter,
	}

	ErrMissingAmountParameter = rest.ValidationError{
		Location: "body.amount",
		Err:      rest.ErrMissingParameter,
	}
)

type Transfer struct {
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
	CreatedAt            string `json:"created_at"`
}

type ListTransfersResponse struct {
	Transfers []Transfer `json:"transfers"`
}

func MapToListTransfersResponse(transfers []entity.Transfer) ListTransfersResponse {
	response := ListTransfersResponse{
		Transfers: make([]Transfer, 0, len(transfers)),
	}

	for _, transfer := range transfers {
		response.Transfers = append(response.Transfers, Transfer{
			AccountOriginID:      transfer.AccountOriginID.String(),
			AccountDestinationID: transfer.AccountDestinationID.String(),
			Amount:               transfer.Amount,
			CreatedAt:            transfer.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return response
}

type PerformTransferInput struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

func (t PerformTransferInput) IsValid() error {
	var errs rest.ValidationErrors

	if t.AccountDestinationID == "" {
		errs = append(errs, ErrMissingAccountDestIDParameter)
	}

	if t.Amount == 0 {
		errs = append(errs, ErrMissingAmountParameter)
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

type PerformTransferResponse struct {
	AccountOriginID      string `json:"account_origin_id"`
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

func MapToPerformTransferResponse(originID, destinationID string, amount int) PerformTransferResponse {
	return PerformTransferResponse{
		AccountOriginID:      originID,
		AccountDestinationID: destinationID,
		Amount:               amount,
	}
}
