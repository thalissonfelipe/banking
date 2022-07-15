package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := make([]*proto.Account, 0)

	for _, acc := range accounts {
		response = append(response, &proto.Account{
			Id:        acc.ID.String(),
			Name:      acc.Name,
			Cpf:       acc.CPF.String(),
			Balance:   int64(acc.Balance),
			CreatedAt: timestamppb.New(acc.CreatedAt),
		})
	}

	return &proto.ListAccountsResponse{Accounts: response}, nil
}

func (h Handler) GetAccountBalance(
	ctx context.Context, req *proto.GetAccountBalanceRequest,
) (*proto.GetAccountBalanceResponse, error) {
	accountID, err := uuid.Parse(req.GetAccountId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account id")
	}

	balance, err := h.accountUsecase.GetAccountBalanceByID(ctx, vos.AccountID(accountID))
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account not found")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.GetAccountBalanceResponse{Balance: int64(balance)}, nil
}

func (h Handler) CreateAccount(
	ctx context.Context, req *proto.CreateAccountRequest,
) (*proto.CreateAccountResponse, error) {
	var errs []*errdetails.BadRequest_FieldViolation

	if req.GetName() == "" {
		errs = append(errs, newFieldViolation("name", "must not be empty"))
	}

	account, err := entity.NewAccount(req.Name, req.Cpf, req.Secret)
	if err != nil {
		if errors.Is(err, vos.ErrInvalidCPF) {
			errs = append(errs, newFieldViolation("cpf", "invalid value"))
		}

		if errors.Is(err, vos.ErrInvalidSecret) {
			errs = append(errs, newFieldViolation("secret", "invalid value"))
		}
	}

	if len(errs) != 0 {
		return nil, newBadRequestError(errs)
	}

	err = h.accountUsecase.CreateAccount(ctx, &account)
	if err != nil {
		if errors.Is(err, entity.ErrAccountAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, "account already exists")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.CreateAccountResponse{Id: account.ID.String()}, nil
}
