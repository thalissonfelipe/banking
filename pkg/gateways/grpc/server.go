package grpc

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/transfer"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	accountUsecase  account.Usecase
	transferUsecase transfer.UseCase
	auth            *auth.Auth
	proto.UnimplementedBankingServiceServer
}

func NewServer(accountUsecase account.Usecase, transferUsecase transfer.UseCase, auth *auth.Auth) *Server {
	return &Server{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		auth:            auth,
	}
}
