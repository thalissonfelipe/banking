package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	proto "github.com/thalissonfelipe/banking/gen/banking/v1"
)

func (h Handler) ListAccounts(ctx context.Context, _ *proto.ListAccountsRequest) (*proto.ListAccountsResponse, error) {
	accounts, err := h.accountUsecase.ListAccounts(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	response := make([]*proto.Account, 0)

	for _, acc := range accounts {
		response = append(response, domainAccountToGRPC(acc))
	}

	return &proto.ListAccountsResponse{Accounts: response}, nil
}

func (h Handler) GetAccountBalance(
	ctx context.Context, request *proto.GetAccountBalanceRequest) (*proto.GetAccountBalanceResponse, error) {
	accountID, err := uuid.Parse(request.AccountId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account id")
	}

	balance, err := h.accountUsecase.GetAccountBalanceByID(ctx, vos.AccountID(accountID))
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return nil, status.Errorf(codes.NotFound, "account does not exist")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.GetAccountBalanceResponse{Balance: int64(balance)}, nil
}

func (h Handler) CreateAccount(
	ctx context.Context, request *proto.CreateAccountRequest) (*proto.CreateAccountResponse, error) {
	// TODO: refact
	if request.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing name parameter")
	}

	if request.Cpf == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing cpf parameter")
	}

	if request.Secret == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing secret parameter")
	}

	account, err := entity.NewAccount(request.Name, request.Cpf, request.Secret)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request parameters")
	}

	err = h.accountUsecase.CreateAccount(ctx, &account)
	if err != nil {
		if errors.Is(err, entity.ErrAccountAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "account already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.CreateAccountResponse{Id: account.ID.String()}, nil
}

func domainAccountToGRPC(account entity.Account) *proto.Account {
	return &proto.Account{
		Id:        account.ID.String(),
		Name:      account.Name,
		Cpf:       account.CPF.String(),
		Balance:   int64(account.Balance),
		CreatedAt: timestamppb.New(account.CreatedAt),
	}
}
