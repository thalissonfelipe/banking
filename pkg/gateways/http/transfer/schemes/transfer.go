package schemes

import "github.com/thalissonfelipe/banking/pkg/gateways/http/responses"

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

func (t CreateTransferInput) IsValid() error {
	if t.AccountDestinationID == "" {
		return responses.ErrMissingAccDestinationIDParameter
	}

	if t.Amount == 0 {
		return responses.ErrMissingAmountParameter
	}

	return nil
}
