package schema

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entities"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

func TestSchema_MapToListTransferResponse(t *testing.T) {
	transfer, err := entities.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 100, 100)
	require.NoError(t, err)

	transfer.CreatedAt = time.Now()

	tests := []struct {
		name      string
		transfers []entities.Transfer
		want      ListTransfersResponse
	}{
		{
			name:      "empty list of transfers",
			transfers: nil,
			want:      ListTransfersResponse{Transfers: []Transfer{}},
		},
		{
			name:      "map transfers successfully",
			transfers: []entities.Transfer{transfer, transfer, transfer, transfer},
			want: ListTransfersResponse{Transfers: []Transfer{
				{
					AccountOriginID:      transfer.AccountOriginID.String(),
					AccountDestinationID: transfer.AccountDestinationID.String(),
					Amount:               transfer.Amount,
					CreatedAt:            transfer.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					AccountOriginID:      transfer.AccountOriginID.String(),
					AccountDestinationID: transfer.AccountDestinationID.String(),
					Amount:               transfer.Amount,
					CreatedAt:            transfer.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					AccountOriginID:      transfer.AccountOriginID.String(),
					AccountDestinationID: transfer.AccountDestinationID.String(),
					Amount:               transfer.Amount,
					CreatedAt:            transfer.CreatedAt.UTC().Format(time.RFC3339),
				},
				{
					AccountOriginID:      transfer.AccountOriginID.String(),
					AccountDestinationID: transfer.AccountDestinationID.String(),
					Amount:               transfer.Amount,
					CreatedAt:            transfer.CreatedAt.UTC().Format(time.RFC3339),
				},
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MapToListTransfersResponse(tt.transfers)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_PerformTransferInput_IsValid(t *testing.T) {
	tests := []struct {
		name    string
		input   PerformTransferInput
		wantErr error
	}{
		{
			name: "should validate input without errors",
			input: PerformTransferInput{
				AccountDestinationID: vos.NewAccountID().String(),
				Amount:               100,
			},
			wantErr: nil,
		},
		{
			name: "should return err acc destination id is blank",
			input: PerformTransferInput{
				Amount: 100,
			},
			wantErr: rest.ErrMissingAccDestinationIDParameter,
		},
		{
			name: "should return err if amount is blank",
			input: PerformTransferInput{
				AccountDestinationID: vos.NewAccountID().String(),
			},
			wantErr: rest.ErrMissingAmountParameter,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.IsValid()
			assert.ErrorIs(t, got, tt.wantErr)
		})
	}
}

func TestSchema_MapToPerformTransferResponse(t *testing.T) {
	originID := vos.NewAccountID().String()
	destinationID := vos.NewAccountID().String()
	amount := 100

	want := PerformTransferResponse{
		AccountOriginID:      originID,
		AccountDestinationID: destinationID,
		Amount:               amount,
	}

	got := MapToPerformTransferResponse(originID, destinationID, amount)
	assert.Equal(t, want, got)
}
