package grpc

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	usecase account.Usecase
	auth    *auth.Auth
	proto.UnimplementedBankingServiceServer
}

func NewServer(usecase account.Usecase, auth *auth.Auth) *Server {
	return &Server{
		usecase: usecase,
		auth:    auth,
	}
}
