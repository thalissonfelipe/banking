package grpc

import (
	"github.com/thalissonfelipe/banking/banking/domain/usecases"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	accountUsecase  usecases.Account
	transferUsecase usecases.Transfer
	authUsecase     usecases.Auth
}

func NewHandler(accountUsecase usecases.Account, transferUsecase usecases.Transfer, authUsecase usecases.Auth) *Handler {
	return &Handler{
		accountUsecase:  accountUsecase,
		transferUsecase: transferUsecase,
		authUsecase:     authUsecase,
	}
}

func newFieldViolation(field, desc string) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: desc,
	}
}

func newBadRequestError(errs []*errdetails.BadRequest_FieldViolation) error {
	st := status.New(codes.InvalidArgument, "invalid parameters")
	br := &errdetails.BadRequest{FieldViolations: errs}

	st, err := st.WithDetails(br)
	if err != nil {
		return status.Error(codes.Internal, "internal server error")
	}

	return st.Err()
}
