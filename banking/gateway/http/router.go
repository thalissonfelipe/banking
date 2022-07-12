package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	httpSwagger "github.com/swaggo/http-swagger"

	acc_usecase "github.com/thalissonfelipe/banking/banking/domain/account/usecase"
	tr_usecase "github.com/thalissonfelipe/banking/banking/domain/transfer/usecase"
	acc_repo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	tr_repo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/hash"
	acc_handler "github.com/thalissonfelipe/banking/banking/gateway/http/account"
	auth_handler "github.com/thalissonfelipe/banking/banking/gateway/http/auth"
	tr_handler "github.com/thalissonfelipe/banking/banking/gateway/http/transfer"
	"github.com/thalissonfelipe/banking/banking/services/auth"
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
	const timeout = 60

	router.Use(middleware.Timeout(timeout * time.Second))

	router.Route("/api/v1", func(r chi.Router) {
		acc_handler.NewHandler(r, accountUsecase)
		auth_handler.NewHandler(r, authService)
		tr_handler.NewHandler(r, transferUsecase)
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"),
	))

	return router
}
