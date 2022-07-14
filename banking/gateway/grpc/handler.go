package grpc

import "github.com/thalissonfelipe/banking/banking/domain/usecases"

type Handler struct {
	accountUsecase  usecases.Account
	transferUsecase usecases.Transfer
	authUsecase     usecases.Auth
}

func NewHandler(accountUsecase usecases.Account, transferUsecase usecases.Transfer, authUsecase usecases.Auth) *Handler {
	return &Handler{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		authUsecase:     authUsecase,
	}
}
