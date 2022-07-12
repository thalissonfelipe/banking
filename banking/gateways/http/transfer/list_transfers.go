package transfer

import (
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateways/http/rest"
	"github.com/thalissonfelipe/banking/banking/gateways/http/transfer/schemes"
	"github.com/thalissonfelipe/banking/banking/services/auth"
)

// ListTransfers returns all transfers
// @Tags transfers
// @Summary List transfers
// @Description List all transfers.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Success 200 {array} transferResponse
// @Failure 401 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /transfers [GET].
func (h Handler) ListTransfers(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(r.Header.Get("Authorization"))
	accountID := vos.ConvertStringToAccountID(auth.GetIDFromToken(token))

	transfers, err := h.usecase.ListTransfers(r.Context(), accountID)
	if err != nil {
		rest.SendError(w, http.StatusInternalServerError, rest.ErrInternalError)

		return
	}

	transfersResponse := make([]schemes.TransferListResponse, 0)

	for _, transfer := range transfers {
		transfersResponse = append(transfersResponse, convertTransferToTransferListResponse(transfer))
	}

	rest.SendJSON(w, http.StatusOK, transfersResponse)
}
