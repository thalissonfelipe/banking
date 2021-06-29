package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/thalissonfelipe/banking/pkg/domain/entities"
	"github.com/thalissonfelipe/banking/pkg/domain/vos"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

func (s Server) GetTransfers(ctx context.Context, _ *proto.ListTransfersRequest) (*proto.ListTransfersResponse, error) {
	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing context metadata")
	}

	token := meta["authorization"][0]

	accountID, err := uuid.Parse(auth.GetIDFromToken(token))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid account id")
	}

	transfers, err := s.transferUsecase.ListTransfers(ctx, vos.AccountID(accountID))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server error")
	}

	var response []*proto.Transfer

	for _, tr := range transfers {
		response = append(response, domainTransferToGRPC(tr))
	}

	return &proto.ListTransfersResponse{Transfers: response}, nil
}

func domainTransferToGRPC(transfer entities.Transfer) *proto.Transfer {
	return &proto.Transfer{
		Id:                   transfer.ID.String(),
		AccountDestinationId: transfer.AccountDestinationID.String(),
		Amount:               int64(transfer.Amount),
		CreatedAt:            timestamppb.New(transfer.CreatedAt),
	}
}
