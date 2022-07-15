package grpc

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
	proto "github.com/thalissonfelipe/banking/gen/banking/v1"
)

func (h Handler) ListTransfers(
	ctx context.Context, _ *proto.ListTransfersRequest,
) (*proto.ListTransfersResponse, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing context metadata")
	}

	token := meta["authorization"][0]

	accountID, err := uuid.Parse(jwt.GetAccountIDFromToken(token))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account id")
	}

	transfers, err := h.transferUsecase.ListTransfers(ctx, vos.AccountID(accountID))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	response := make([]*proto.Transfer, 0)

	for _, transfer := range transfers {
		response = append(response, &proto.Transfer{
			Id:                   transfer.ID.String(),
			AccountDestinationId: transfer.AccountDestinationID.String(),
			Amount:               int64(transfer.Amount),
			CreatedAt:            timestamppb.New(transfer.CreatedAt),
		})
	}

	return &proto.ListTransfersResponse{Transfers: response}, nil
}

func (h Handler) PerformTransfer(
	ctx context.Context, request *proto.PerformTransferRequest,
) (*proto.PerformTransferResponse, error) {
	var errs []*errdetails.BadRequest_FieldViolation

	if request.GetAccountDestinationId() == "" {
		errs = append(errs, newFieldViolation("account_origin_id", "must not be empty"))
	}

	if request.Amount <= 0 {
		errs = append(errs, newFieldViolation("amount", "must be a value bigger than 0"))
	}

	if len(errs) != 0 {
		return nil, newBadRequestError(errs)
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing context metadata")
	}

	token := meta["authorization"][0]

	accounOriginID, err := uuid.Parse(jwt.GetAccountIDFromToken(token))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account origin id")
	}

	accounDestinationID, err := uuid.Parse(request.AccountDestinationId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account destination id")
	}

	if accounOriginID == accounDestinationID {
		return nil, status.Error(
			codes.InvalidArgument,
			"account origin id cannot be equal to destination id")
	}

	input := usecases.NewPerformTransferInput(
		vos.AccountID(accounOriginID),
		vos.AccountID(accounDestinationID),
		int(request.Amount),
	)

	err = h.transferUsecase.PerformTransfer(ctx, input)
	if err != nil {
		if errors.Is(err, entity.ErrAccountNotFound) {
			return nil, status.Error(codes.NotFound, "account origin does not exist")
		}

		if errors.Is(err, entity.ErrAccountDestinationNotFound) {
			return nil, status.Error(codes.NotFound, "account destination does not exist")
		}

		if errors.Is(err, entity.ErrInsufficientFunds) {
			return nil, status.Error(codes.InvalidArgument, "insufficient funds")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.PerformTransferResponse{}, nil
}
