package schema

import "github.com/thalissonfelipe/banking/banking/gateway/http/rest"

var (
	ErrMissingCPFParameter = rest.ValidationError{
		Location: "body.cpf",
		Err:      rest.ErrMissingParameter,
	}

	ErrMissingSecretParameter = rest.ValidationError{
		Location: "body.secret",
		Err:      rest.ErrMissingParameter,
	}
)

type LoginResponse struct {
	Token string `json:"token"`
}

func MapToLoginResponse(token string) LoginResponse {
	return LoginResponse{Token: token}
}

type LoginInput struct {
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

func (r LoginInput) IsValid() error {
	var errs rest.ValidationErrors

	if r.CPF == "" {
		errs = append(errs, ErrMissingCPFParameter)
	}

	if r.Secret == "" {
		errs = append(errs, ErrMissingSecretParameter)
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}
