package grpc

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/services"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	accountUsecase  usecases.Account
	transferUsecase usecases.Transfer
	auth            services.Auth
	proto.UnimplementedBankingServiceServer
}

func NewServer(accountUsecase usecases.Account, transferUsecase usecases.Transfer, auth services.Auth) *Server {
	return &Server{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		auth:            auth,
	}
}
