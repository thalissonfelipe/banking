package transfer

import (
	"net/http"
	"strings"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateway/http/transfer/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

// ListTransfers returns all transfers.
// @Tags Transfers
// @Summary Lists all transfers.
// @Description Lists all transfers. User must be authenticated.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Success 200 {object} schema.ListTransfersResponse
// @Failure 401 {object} rest.UnauthorizedError
// @Failure 500 {object} rest.InternalServerError
// @Router /transfers [GET].
func (h Handler) ListTransfers(r *http.Request) rest.Response {
	token := strings.Split(r.Header.Get("Authorization"), "Bearer ")[1]

	accountID, err := rest.ParseUUID(jwt.GetAccountIDFromToken(token), "token")
	if err != nil {
		return rest.BadRequest(err, "invalid token")
	}

	transfers, err := h.usecase.ListTransfers(r.Context(), vos.AccountID(accountID))
	if err != nil {
		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToListTransfersResponse(transfers))
}
