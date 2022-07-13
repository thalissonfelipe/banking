package grpc

import (
	"github.com/thalissonfelipe/banking/banking/domain/account"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/services"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	accountUsecase  account.Usecase
	transferUsecase transfer.Usecase
	auth            services.Auth
	proto.UnimplementedBankingServiceServer
}

func NewServer(accountUsecase account.Usecase, transferUsecase transfer.Usecase, auth services.Auth) *Server {
	return &Server{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		auth:            auth,
	}
}
