package transfer

import (
	"net/http"
	"strings"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/services/auth"
)

// ListTransfers returns all transfers
// @Tags transfers
// @Summary List transfers
// @Description List all transfers.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Success 200 schema.ListTransfersResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /transfers [GET].
func (h Handler) ListTransfers(r *http.Request) rest.Response {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	accountID, err := rest.ParseUUID(auth.GetIDFromToken(token), "token")
	if err != nil {
		return rest.BadRequest(err, "invalid token")
	}

	transfers, err := h.usecase.ListTransfers(r.Context(), vos.AccountID(accountID))
	if err != nil {
		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToListTransfersResponse(transfers))
}
