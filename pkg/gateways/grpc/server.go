package grpc

import (
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type Server struct {
	usecase account.Usecase
	proto.UnimplementedBankingServiceServer
}

func NewServer(usecase account.Usecase) *Server {
	return &Server{
		usecase: usecase,
	}
}
