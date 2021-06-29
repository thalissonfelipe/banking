package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/pkg/domain/account"
	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

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

func (s Server) CreateAccount(ctx context.Context, request *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	// TODO: refact
	if request.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "missing name parameter")
	}

	if request.Cpf == "" {
		return nil, status.Error(codes.InvalidArgument, "missing cpf parameter")
	}

	if request.Secret == "" {
		return nil, status.Error(codes.InvalidArgument, "missing secret parameter")
	}

	cpf, err := vos.NewCPF(request.Cpf)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	secret, err := vos.NewSecret(request.Secret)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	input := account.NewCreateAccountInput(request.Name, cpf, secret)

	account, err := s.usecase.CreateAccount(ctx, input)
	if err != nil {
		if errors.Is(err, entities.ErrAccountAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "account already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.CreateAccountResponse{Id: account.ID.String()}, nil
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
