package transfer

import (
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

// ListTransfers returns all transfers
// @Tags transfers
// @Summary List transfers
// @Description List all transfers.
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer Authorization Token"
// @Success 200 {array} transferResponse
// @Failure 401 {string} string "No token provided"
// @Failure 500 {string} string "Internal server error"
// @Router /transfers [GET]
func (h Handler) ListTransfers(w http.ResponseWriter, r *http.Request) {
	token := getTokenFromHeader(r.Header.Get("Authorization"))
	accountID := vos.ConvertStringToID(auth.GetIDFromToken(token))
	transfers, err := h.usecase.ListTransfers(r.Context(), accountID)
	if err != nil {
		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	transfersResponse := make([]transferResponse, 0)
	for _, transfer := range transfers {
		transfersResponse = append(transfersResponse, convertTransferToTransferResponse(transfer))
	}

	responses.SendJSON(w, http.StatusOK, transfersResponse)
}
