package rest

type InvalidCredentialsError struct {
	Error string `json:"error" example:"invalid credentials"`
}

type UnauthorizedError struct {
	Error string `json:"error" example:"unauthorized"`
}

type NotFoundError struct {
	Error string `json:"error" example:"account not found"`
}

type ConflictError struct {
	Error string `json:"error" example:"account already exists"`
}

type InternalServerError struct {
	Error string `json:"error" example:"internal server error"`
}
