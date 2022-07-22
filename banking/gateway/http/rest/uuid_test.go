package rest

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestParseUUID(t *testing.T) {
	t.Parallel()

	wantID := uuid.New()

	tests := []struct {
		name    string
		id      string
		loc     string
		wantID  uuid.UUID
		wantErr error
	}{
		{
			name:    "should parse uuid successfully",
			id:      wantID.String(),
			loc:     "",
			wantID:  wantID,
			wantErr: nil,
		},
		{
			name:    "should return an error if uuid is invalid",
			id:      "invalid",
			loc:     "path",
			wantID:  uuid.UUID{},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseUUID(tt.id, tt.loc)
			if err != nil {
				var verr ValidationError
				assert.True(t, errors.As(err, &verr))
				assert.Equal(t, tt.loc, verr.Location)
				assert.ErrorIs(t, ErrInvalidUUID, verr.Err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.wantID, got)
		})
	}
}
