package grpc

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	accountUsecase  usecases.Account
	transferUsecase usecases.Transfer
	authUsecase     usecases.Auth
	proto.UnimplementedBankingServiceServer
}

func NewServer(accountUsecase usecases.Account, transferUsecase usecases.Transfer, authUsecase usecases.Auth) *Server {
	return &Server{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		authUsecase:     authUsecase,
	}
}
