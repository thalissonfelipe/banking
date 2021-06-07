package account

import (
	"encoding/json"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
)

// CreateAccount creates a new account
// @Tags accounts
// @Summary Creates a new account
// @Description Creates a new account given a name, cpf and secret.
// @Description Secret must be a minimum of 8, a maximum of 20, at least one lowercase character,
// @Description one uppercase character and one number.
// @Description CPF must have the format XXX.XXX.XXX-XX or XXXXXXXXXXX.
// @Accept json
// @Produce json
// @Param Body body requestBody true "Body"
// @Success 201 {object} createdAccountResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 409 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /accounts [POST]
func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body requestBody

	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, responses.ErrInvalidJSON)

		return
	}

	if err := body.isValid(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err)

		return
	}

	cpf, err := vos.NewCPF(body.CPF)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)

		return
	}

	secret, err := vos.NewSecret(body.Secret)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err)

		return
	}

	input := account.NewCreateAccountInput(body.Name, cpf, secret)

	acc, err := h.usecase.CreateAccount(r.Context(), input)
	if err != nil {
		responses.HandleError(w, err)

		return
	}

	accResponse := convertAccountToCreatedAccountResponse(acc)
	responses.SendJSON(w, http.StatusCreated, accResponse)
}
