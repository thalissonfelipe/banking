package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	proto "github.com/thalissonfelipe/banking/gen/banking/v1"
)

func (h Handler) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := h.authUsecase.Autheticate(ctx, request.GetCpf(), request.GetSecret())
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &proto.LoginResponse{Token: token}, nil
}
