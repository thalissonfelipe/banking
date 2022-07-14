package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

func (s Server) Login(ctx context.Context, request *proto.LoginRequest) (*proto.LoginResponse, error) {
	token, err := s.authUsecase.Autheticate(ctx, request.Cpf, request.Secret)
	if err != nil {
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			return nil, status.Errorf(codes.InvalidArgument, "wrong credentials")
		}

		return nil, status.Errorf(codes.Internal, "internal server error")
	}

	return &proto.LoginResponse{Token: token}, nil
}
