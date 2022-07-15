package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

//nolint:gochecknoglobals
var jwtAllowedMethods = map[string]struct{}{
	"/banking.BankingService/ListTransfers":   {},
	"/banking.BankingService/PerformTransfer": {},
}

func AuthInterceptor(
	ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (interface{}, error) {
	_, ok := jwtAllowedMethods[info.FullMethod]
	if !ok {
		return handler(ctx, req)
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing context metadata")
	}

	if len(meta["authorization"]) != 1 {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	token := meta["authorization"][0]

	if err := jwt.IsTokenValid(token); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return handler(ctx, req)
}
