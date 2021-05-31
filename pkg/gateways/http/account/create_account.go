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
