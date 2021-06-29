package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
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

func (s Server) GetAccountBalance(ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	accountID, err := uuid.Parse(request.AccountId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account id")
	}

	balance, err := s.usecase.GetAccountBalanceByID(ctx, vos.AccountID(accountID))
	if err != nil {
		if errors.Is(err, entities.ErrAccountDoesNotExist) {
			return nil, status.Error(codes.NotFound, "account does not exist")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.GetAccountBalanceResponse{Balance: int64(balance)}, nil
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