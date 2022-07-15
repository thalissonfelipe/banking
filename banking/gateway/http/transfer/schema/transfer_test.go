package schema

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

func TestSchema_MapToListTransferResponse(t *testing.T) {
	t.Parallel()

	transfer, err := entity.NewTransfer(vos.NewAccountID(), vos.NewAccountID(), 100, 100)
	require.NoError(t, err)

	transfer.CreatedAt = time.Now()

	tests := []struct {
		name      string
		transfers []entity.Transfer
		want      ListTransfersResponse
	}{
		{
			name:      "empty list of transfers",
			transfers: nil,
			want:      ListTransfersResponse{Transfers: []Transfer{}},
		},
		{
			name:      "map transfers successfully",
			transfers: []entity.Transfer{transfer, transfer, transfer, transfer},
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
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := MapToListTransfersResponse(tt.transfers)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSchema_PerformTransferInput_IsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   PerformTransferInput
		wantErr rest.ValidationErrors
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
			name: "should return err if acc destination id is blank",
			input: PerformTransferInput{
				Amount: 100,
			},
			wantErr: rest.ValidationErrors{ErrMissingAccountDestIDParameter},
		},
		{
			name: "should return err if amount is blank",
			input: PerformTransferInput{
				AccountDestinationID: vos.NewAccountID().String(),
			},
			wantErr: rest.ValidationErrors{ErrMissingAmountParameter},
		},
		{
			name:    "should return err if acc destination id and amount are blank",
			input:   PerformTransferInput{},
			wantErr: rest.ValidationErrors{ErrMissingAccountDestIDParameter, ErrMissingAmountParameter},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := tt.input.IsValid()
			if err != nil {
				var errs rest.ValidationErrors
				require.True(t, errors.As(err, &errs))

				assert.Len(t, errs, len(tt.wantErr))

				for i, e := range errs {
					var verr rest.ValidationError
					require.True(t, errors.As(e, &verr))

					assert.ErrorIs(t, verr, tt.wantErr[i])
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSchema_MapToPerformTransferResponse(t *testing.T) {
	t.Parallel()

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
