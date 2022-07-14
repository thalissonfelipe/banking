package auth

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// Login logs in :D
// @Tags login
// @Summary Log in
// @Description Returns a JWT to be used on /transfers endpoints.
// @Accept json
// @Produce json
// @Param Body body requestBody true "Body"
// @Success 200 {object} schema.LoginResponse
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /login [POST].
func (h Handler) Login(r *http.Request) rest.Response {
	var body schema.LoginInput
	if err := rest.DecodeRequestBody(r, &body); err != nil {
		return rest.BadRequest(err, "invalid request body")
	}

	token, err := h.usecase.Autheticate(r.Context(), body.CPF, body.Secret)
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			return rest.InvalidCredentials(err)
		}

		return rest.InternalServer(err)
	}

	return rest.OK(schema.MapToLoginResponse(token))
}
