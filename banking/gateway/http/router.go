package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/thalissonfelipe/banking/banking/domain/usecases/account"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/transfer"
	accountRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	transferRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/hash"
	accHandler "github.com/thalissonfelipe/banking/banking/gateway/http/account"
	authHandler "github.com/thalissonfelipe/banking/banking/gateway/http/auth"
	trHandler "github.com/thalissonfelipe/banking/banking/gateway/http/transfer"
	"github.com/thalissonfelipe/banking/banking/services/auth"
)

func NewRouter(db *pgx.Conn) http.Handler {
	// Set dependencies
	hash := hash.Hash{}
	accountRepo := accountRepo.NewRepository(db)
	accountUsecase := account.NewAccountUsecase(accountRepo, hash)
	authService := auth.NewAuth(accountUsecase, hash)
	transferRepo := transferRepo.NewRepository(db)
	transferUsecase := transfer.NewTransferUsecase(transferRepo, accountUsecase)

	router := chi.NewRouter()

	// middlewares
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	const timeout = 60

	router.Use(middleware.Timeout(timeout * time.Second))

	router.Route("/api/v1", func(r chi.Router) {
		accHandler.NewHandler(r, accountUsecase)
		authHandler.NewHandler(r, authService)
		trHandler.NewHandler(r, transferUsecase)
	})

	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:5000/swagger/doc.json"),
	))

	return router
}
