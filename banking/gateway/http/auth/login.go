package auth

import (
	"errors"
	"net/http"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/gateway/http/auth/schema"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
)

// Login logs in the server.
// @Tags Sign In
// @Summary Logs in the server.
// @Description Returns a JWT to be used on /transfers endpoints.
// @Accept json
// @Produce json
// @Param Body body schema.LoginInput true "Body"
// @Success 200 {object} schema.LoginResponse
// @Failure 400 {object} rest.InvalidCredentialsError
// @Failure 500 {object} rest.InternalServerError
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
