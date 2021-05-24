package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4"

	acc_usecase "github.com/thalissonfelipe/banking/pkg/domain/account/usecase"
	tr_usecase "github.com/thalissonfelipe/banking/pkg/domain/transfer/usecase"
	acc_repo "github.com/thalissonfelipe/banking/pkg/gateways/db/postgres/account"
	tr_repo "github.com/thalissonfelipe/banking/pkg/gateways/db/postgres/transfer"
	"github.com/thalissonfelipe/banking/pkg/gateways/hash"
	acc_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/account"
	auth_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/auth"
	"github.com/thalissonfelipe/banking/pkg/gateways/http/middlewares"
	tr_handler "github.com/thalissonfelipe/banking/pkg/gateways/http/transfer"
	"github.com/thalissonfelipe/banking/pkg/services/auth"
)

func NewRouter(db *pgx.Conn) http.Handler {
	// Set dependencies
	hash := hash.Hash{}
	accountRepo := acc_repo.NewRepository(db)
	accountUsecase := acc_usecase.NewAccountUseCase(accountRepo, hash)
	authService := auth.NewAuth(accountUsecase, hash)
	transferRepo := tr_repo.NewRepository(db)
	transferUsecase := tr_usecase.NewTransfer(transferRepo, accountUsecase)

	// Register endpoints
	router := mux.NewRouter()
	router = router.PathPrefix("/api/v1").Subrouter()

	acc_handler.NewHandler(router, accountUsecase)
	auth_handler.NewHandler(router, authService)
	tr_handler.NewHandler(router, transferUsecase)

	router.Use(middlewares.LogMiddleware)

	return router
}
