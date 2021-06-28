package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
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

func (s Server) GetAccounts(ctx context.Context, _ *proto.ListAccountsRequest) (*proto.ListAccountsResponse, error) {
	accounts, err := s.usecase.ListAccounts(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var response []*proto.Account

	for _, acc := range accounts {
		response = append(response, domainToGRPC(acc))
	}

	return &proto.ListAccountsResponse{Accounts: response}, nil
}

func domainToGRPC(account entities.Account) *proto.Account {
	return &proto.Account{
		Id:        account.ID.String(),
		Name:      account.Name,
		Cpf:       account.CPF.String(),
		Balance:   int64(account.Balance),
		CreatedAt: timestamppb.New(account.CreatedAt),
	}
}
