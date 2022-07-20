package grpc

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/thalissonfelipe/banking/banking/domain/usecases/account"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/auth"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/transfer"
	accountRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	transferRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/hash"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
	proto "github.com/thalissonfelipe/banking/gen/banking/v1"
)

func NewServer(logger *zap.Logger, db *pgxpool.Pool) *grpc.Server {
	logger = logger.With(zap.String("module", "grpc"))

	h, j := hash.New(), jwt.New()

	accRepository := accountRepo.NewRepository(db)
	accountUsecase := account.NewAccountUsecase(accRepository, h)
	trRepository := transferRepo.NewRepository(db)
	transferUsecase := transfer.NewTransferUsecase(trRepository, accountUsecase)
	authUsecase := auth.NewAuth(accountUsecase, h, j)
	handler := NewHandler(accountUsecase, transferUsecase, authUsecase)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(
			grpczap.UnaryServerInterceptor(logger),
			AuthInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		)),
	)

	proto.RegisterBankingServiceServer(server, handler)

	return server
}
