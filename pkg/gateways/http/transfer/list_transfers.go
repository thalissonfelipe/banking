package transfer

import (
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

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
