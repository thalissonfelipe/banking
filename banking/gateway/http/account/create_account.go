package account

import (
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
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
// @Router /accounts [POST].
func (h Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body schema.CreateAccountInput

	if err := rest.DecodeRequestBody(r, &body); err != nil {
		rest.HandleBadRequestError(w, err)

		return
	}

	account, err := entity.NewAccount(body.Name, body.CPF, body.Secret)
	if err != nil {
		rest.SendError(w, http.StatusBadRequest, err)

		return
	}

	err = h.usecase.CreateAccount(r.Context(), &account)
	if err != nil {
		rest.HandleError(w, err)

		return
	}

	rest.SendJSON(w, http.StatusCreated, schema.MapToCreateAccountResponse(account))
}
