package auth

import (
	"net/http"

	"github.com/thalissonfelipe/banking/pkg/gateways/http/auth/schemes"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/responses"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/rest"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

// Login logs in :D
// @Tags login
// @Summary Log in
// @Description Returns a JWT to be used on /transfers endpoints.
// @Accept json
// @Produce json
// @Param Body body requestBody true "Body"
// @Success 200 {object} responseBody
// @Failure 400 {object} responses.ErrorResponse
// @Failure 404 {object} responses.ErrorResponse
// @Failure 500 {object} responses.ErrorResponse
// @Router /login [POST].
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	var body schemes.LoginInput

	err := rest.DecodeRequestBody(r, &body)
	if err != nil {
		responses.HandleBadRequestError(w, err)

		return
	}

	input := auth.AuthenticateInput{
		CPF:    body.CPF,
		Secret: body.Secret,
	}

	token, err := h.authService.Autheticate(r.Context(), input)
	if err != nil {
		responses.HandleError(w, err)

		return
	}

	response := schemes.LoginResponse{Token: token}
	responses.SendJSON(w, http.StatusOK, response)
}
