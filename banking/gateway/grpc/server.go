package grpc

import "github.com/thalissonfelipe/banking/banking/domain/usecases"

type Server struct {
	accountUsecase  usecases.Account
	transferUsecase usecases.Transfer
	authUsecase     usecases.Auth
}

func NewServer(accountUsecase usecases.Account, transferUsecase usecases.Transfer, authUsecase usecases.Auth) *Server {
	return &Server{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		authUsecase:     authUsecase,
	}
}
