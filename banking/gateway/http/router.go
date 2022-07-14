package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v4"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/thalissonfelipe/banking/banking/domain/usecases/account"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/auth"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/transfer"
	accountRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/account"
	transferRepo "github.com/thalissonfelipe/banking/banking/gateway/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/hash"
	accountHandler "github.com/thalissonfelipe/banking/banking/gateway/http/account"
	authHandler "github.com/thalissonfelipe/banking/banking/gateway/http/auth"
	"github.com/thalissonfelipe/banking/banking/gateway/http/middlewares"
	"github.com/thalissonfelipe/banking/banking/gateway/http/rest"
	transferHandler "github.com/thalissonfelipe/banking/banking/gateway/http/transfer"
	"github.com/thalissonfelipe/banking/banking/gateway/jwt"
)

func NewRouter(logger *zap.Logger, db *pgx.Conn) http.Handler {
	r := chi.NewRouter()
	logger = logger.With(zap.String("module", "http-server"))

	const timeout = 60 * time.Second

	r.Use(
		middlewares.Logger(logger),
		middleware.Recoverer,
		middleware.StripSlashes,
		middleware.Timeout(timeout),
	)

	r.Route("/docs/", func(r chi.Router) {
		r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "swagger/index.html", http.StatusMovedPermanently)
		})
		r.Get("/swagger/*", httpSwagger.Handler())
	})

	r.Mount("/api/v1", apiRouter(db))

	return r
}

func apiRouter(db *pgx.Conn) chi.Router {
	r := chi.NewRouter()

	r.Use(
		middlewares.RequestID,
		middlewares.RequestIDToLogger,
	)

	hash := hash.New()
	jwt := jwt.New()

	accRepository := accountRepo.NewRepository(db)
	accountUsecase := account.NewAccountUsecase(accRepository, hash)
	trRepository := transferRepo.NewRepository(db)
	transferUsecase := transfer.NewTransferUsecase(trRepository, accountUsecase)
	authUsecase := auth.NewAuth(accountUsecase, hash, jwt)
	accHandler := accountHandler.NewHandler(accountUsecase)
	trHandler := transferHandler.NewHandler(transferUsecase)
	authHandler := authHandler.NewHandler(authUsecase)

	r.Route("/accounts", func(r chi.Router) {
		r.Get("/", rest.Wrap(accHandler.ListAccounts))
		r.Post("/", rest.Wrap(accHandler.CreateAccount))
		r.Get("/{accountID}/balance", rest.Wrap(accHandler.GetAccountBalance))
	})

	r.Route("/transfers", func(r chi.Router) {
		r.Use(middlewares.Authorize)
		r.Get("/", rest.Wrap(trHandler.ListTransfers))
		r.Post("/", rest.Wrap(trHandler.PerformTransfer))
	})

	r.Post("/login", rest.Wrap(authHandler.Login))

	return r
}
