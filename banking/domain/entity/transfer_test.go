package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
)

func TestNewTransfer(t *testing.T) {
	const amount = 100

	accOriginID := vos.NewAccountID()
	accDestID := vos.NewAccountID()

	tests := []struct {
		name    string
		balance int
		want    Transfer
		wantErr error
	}{
		{
			name:    "should create a transfer successfully",
			balance: 100,
			want: Transfer{
				AccountOriginID:      accOriginID,
				AccountDestinationID: accDestID,
				Amount:               amount,
			},
			wantErr: nil,
		},
		{
			name:    "should return an error if amount is bigger than balance",
			balance: 50,
			want:    Transfer{},
			wantErr: ErrInsufficientFunds,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			transfer, err := NewTransfer(accOriginID, accDestID, amount, tt.balance)
			assert.ErrorIs(t, err, tt.wantErr)

			assert.Equal(t, tt.want.AccountOriginID, transfer.AccountOriginID)
			assert.Equal(t, tt.want.AccountDestinationID, transfer.AccountDestinationID)
			assert.Equal(t, tt.want.Amount, transfer.Amount)
		})
	}
}
