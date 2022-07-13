package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/services/auth"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

func (s Server) GetTransfers(
	ctx context.Context, _ *proto.ListTransfersRequest) (*proto.ListTransfersResponse, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	token := meta["authorization"][0]

	accountID, err := uuid.Parse(auth.GetIDFromToken(token))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account id")
	}

	transfers, err := s.transferUsecase.ListTransfers(ctx, vos.AccountID(accountID))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	response := make([]*proto.Transfer, 0)

	for _, tr := range transfers {
		response = append(response, domainTransferToGRPC(tr))
	}

	return &proto.ListTransfersResponse{Transfers: response}, nil
}

func (s Server) CreateTransfer(
	ctx context.Context, request *proto.CreateTransferRequest) (*proto.CreateTransferResponse, error) {
	if request.AccountDestinationId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "missing account destination id parameter")
	}

	if request.Amount == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "missing amount parameter")
	}

	if request.Amount < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "amount must be bigger than 0")
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	token := meta["authorization"][0]

	accounOriginID, err := uuid.Parse(auth.GetIDFromToken(token))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account origin id")
	}

	accounDestinationID, err := uuid.Parse(request.AccountDestinationId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid account destination id")
	}

	if accounOriginID == accounDestinationID {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"account destination id must be different from account origin id")
	}

	input := transfer.NewPerformTransferInput(
		vos.AccountID(accounOriginID),
		vos.AccountID(accounDestinationID),
		int(request.Amount),
	)

	err = s.transferUsecase.PerformTransfer(ctx, input)
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return nil, status.Errorf(codes.NotFound, "account origin does not exist")
		}

		if errors.Is(err, entity.ErrAccountDestinationNotFound) {
			return nil, status.Errorf(codes.NotFound, "account destination does not exist")
		}

		if errors.Is(err, entity.ErrInsufficientFunds) {
			return nil, status.Errorf(codes.InvalidArgument, "insufficient funds")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.CreateTransferResponse{}, nil
}

func domainTransferToGRPC(transfer entity.Transfer) *proto.Transfer {
	return &proto.Transfer{
		Id:                   transfer.ID.String(),
		AccountDestinationId: transfer.AccountDestinationID.String(),
		Amount:               int64(transfer.Amount),
		CreatedAt:            timestamppb.New(transfer.CreatedAt),
	}
}
