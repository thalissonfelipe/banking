package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"

	acc_usecase "github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	tr_usecase "github.com/thalissonfelipe/banking/pkg/domain/transfer/usecase"
	acc_repo "github.com/thalissonfelipe/banking/pkg/gateways/db/postgres/account"
	tr_repo "github.com/thalissonfelipe/banking/pkg/gateways/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/pkg/gateways/hash"
	acc_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/account"
	auth_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/auth"
	tr_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/transfer"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

func NewRouter(db *pgx.Conn) http.Handler {
	// Set dependencies
	hash := hash.Hash{}
	accountRepo := acc_repo.NewRepository(db)
	accountUsecase := acc_usecase.NewAccountUsecase(accountRepo, hash)
	authService := auth.NewAuth(accountUsecase, hash)
	transferRepo := tr_repo.NewRepository(db)
	transferUsecase := tr_usecase.NewTransferUsecase(transferRepo, accountUsecase)

	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	router.Use(middleware.Timeout(60 * time.Second))

	acc_handler.NewHandler(router, accountUsecase)
	auth_handler.NewHandler(router, authService)
	tr_handler.NewHandler(router, transferUsecase)

	return router
}
