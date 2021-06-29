package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

// jwtMethods is an array of methods that need authentication validation.
var jwtMethods = []string{
	"/banking.BankingService/GetTransfers",
	"/banking.BankingService/CreateTransfer",
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ok := needAuthentication(info.FullMethod)
	if !ok {
		return handler(ctx, req)
	}

	meta, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing context metadata")
	}

	if len(meta["authorization"]) != 1 {
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	token := meta["authorization"][0]
	err := auth.IsValidToken(token)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "incorrect access token")
	}

	return handler(ctx, req)
}

// needAuthentication will check if the provided method need authentication validation. Only transfers
// methods will need jwt authentication.
func needAuthentication(method string) bool {
	for _, m := range jwtMethods {
		if m == method {
			return true
		}
	}

	return false
}
