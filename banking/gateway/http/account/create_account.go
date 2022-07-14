package account

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/gateway/http/account/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// CreateAccount creates a new account.
// @Tags Accounts
// @Summary Creates a new account.
// @Description Creates a new account given a name, cpf and secret.
// @Description Secret must be a minimum of 8, a maximum of 20, at least one lowercase character,
// @Description one uppercase character and one number.
// @Description CPF must have the format XXX.XXX.XXX-XX or XXXXXXXXXXX.
// @Accept json
// @Produce json
// @Param Body body schema.CreateAccountInput true "Body"
// @Success 201 {object} schema.CreateAccountResponse
// @Failure 400 {object} rest.Error
// @Failure 409 {object} rest.ConflictError
// @Failure 500 {object} rest.InternalServerError
// @Router /accounts [POST].
func (h Handler) CreateAccount(r *http.Request) rest.Response {
	var body schema.CreateAccountInput
	if err := rest.DecodeRequestBody(r, &body); err != nil {
		return rest.BadRequest(err, "invalid request body")
	}

	account, err := entity.NewAccount(body.Name, body.CPF, body.Secret)
	if err != nil {
		return rest.BadRequest(err, "invalid request body")
	}

	err = h.usecase.CreateAccount(r.Context(), &account)
	if err != nil {
		if errors.Is(err, entity.ErrAccountAlreadyExists) {
			return rest.Conflict(err, "account already exists")
		}

		return rest.InternalServer(err)
	}

	return rest.Created(schema.MapToCreateAccountResponse(account))
}
