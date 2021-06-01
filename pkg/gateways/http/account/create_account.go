package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
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
// @Failure 400 {string} string "Bad request"
// @Failure 409 {string} string "Account already exists"
// @Failure 500 {string} string "Internal server error"
// @Router /accounts [POST]
func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body requestBody
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, errInvalidJSON.Error())
		return
	}

	if err := body.isValid(); err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	cpf, err := vos.NewCPF(body.CPF)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	secret, err := vos.NewSecret(body.Secret)
	if err != nil {
		responses.SendError(w, http.StatusBadRequest, err.Error())
		return
	}

	input := account.NewCreateAccountInput(body.Name, cpf, secret)
	acc, err := h.usecase.CreateAccount(r.Context(), input)
	if err != nil {
		if errors.Is(err, entities.ErrAccountAlreadyExists) {
			responses.SendError(w, http.StatusConflict, err.Error())
			return
		}

		responses.SendError(w, http.StatusInternalServerError, err.Error())
		return
	}

	accResponse := convertAccountToCreatedAccountResponse(acc)
	responses.SendJSON(w, http.StatusCreated, accResponse)
}
