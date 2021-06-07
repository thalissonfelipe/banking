package schemes

import "errors"

var (
	ErrMissingAccDestIDParameter = errors.New("missing account destination id parameter")
	ErrMissingAmountParameter    = errors.New("missing amount parameter")
	ErrDestIDEqualCurrentID      = errors.New("account destination cannot be the account origin id")
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

func (t CreateTransferInput) IsValid() error {
	if t.AccountDestinationID == "" {
		return ErrMissingAccDestIDParameter
	}
	if t.Amount == 0 {
		return ErrMissingAmountParameter
	}

	return nil
}
