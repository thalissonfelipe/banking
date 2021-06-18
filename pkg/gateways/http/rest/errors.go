package rest

import "errors"

// REST API generic errors.
var (
	ErrInvalidJSON                      = errors.New("invalid json")
	ErrAccountNotFound                  = errors.New("account does not exist")
	ErrAccountOriginNotFound            = errors.New("account origin does not exist")
	ErrAccountDestinationNotFound       = errors.New("account destination does not exist")
	ErrInternalError                    = errors.New("internal server error")
	ErrInsufficientFunds                = errors.New("insufficient funds")
	ErrInvalidCredentials               = errors.New("cpf or secret are invalid")
	ErrAccountAlreadyExists             = errors.New("account already exists")
	ErrMissingNameParameter             = errors.New("missing name parameter")
	ErrMissingCPFParameter              = errors.New("missing cpf parameter")
	ErrMissingSecretParameter           = errors.New("missing secret parameter")
	ErrMissingAccDestinationIDParameter = errors.New("missing account destination id parameter")
	ErrMissingAmountParameter           = errors.New("missing amount parameter")
	ErrDestinationIDEqToOriginID        = errors.New("account destination cannot be the account origin id")
)
