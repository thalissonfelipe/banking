package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/thalissonfelipe/banking/banking/domain/entity"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/account"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/auth"
	"github.com/thalissonfelipe/banking/banking/domain/usecases/transfer"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	grpcServer "github.com/thalissonfelipe/banking/banking/gateway/grpc"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

// TODO: remove in memory mocks

type InMemoryAccountDatabase struct {
	Accounts []entity.Account
}

func (i InMemoryAccountDatabase) ListAccounts(ctx context.Context) ([]entity.Account, error) {
	return i.Accounts, nil
}

func (i InMemoryAccountDatabase) GetAccountBalanceByID(ctx context.Context, id vos.AccountID) (int, error) {
	for _, acc := range i.Accounts {
		if acc.ID == id {
			return acc.Balance, nil
		}
	}

	return 0, entity.ErrAccountNotFound
}

func (i *InMemoryAccountDatabase) CreateAccount(ctx context.Context, account *entity.Account) error {
	for _, acc := range i.Accounts {
		if acc.CPF == account.CPF {
			return entity.ErrAccountAlreadyExists
		}
	}

	i.Accounts = append(i.Accounts, *account)

	return nil
}

func (i InMemoryAccountDatabase) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entity.Account, error) {
	for _, acc := range i.Accounts {
		if acc.CPF == cpf {
			return acc, nil
		}
	}

	return entity.Account{}, entity.ErrAccountNotFound
}

func (i InMemoryAccountDatabase) GetAccountByID(ctx context.Context, id vos.AccountID) (entity.Account, error) {
	for _, acc := range i.Accounts {
		if acc.ID == id {
			return acc, nil
		}
	}

	return entity.Account{}, entity.ErrAccountNotFound
}

type InMemoryEncrypter struct{}

func (i InMemoryEncrypter) Hash(secret string) ([]byte, error) {
	return []byte(secret), nil
}

func (i InMemoryEncrypter) CompareHashAndSecret(hashedSecret, secret []byte) error {
	return nil
}

type InMemoryJWT struct{}

func (i InMemoryJWT) NewToken(accountID string) (string, error) {
	return "token", nil
}

type InMemoryTransferDatabase struct {
	Transfers []entity.Transfer
}

func (i InMemoryTransferDatabase) ListTransfers(ctx context.Context, id vos.AccountID) ([]entity.Transfer, error) {
	return i.Transfers, nil
}

func (i *InMemoryTransferDatabase) PerformTransfer(ctx context.Context, transfer *entity.Transfer) error {
	return nil
}

func main() {
	creds, err := credentials.NewServerTLSFromFile("../cert/server.crt", "../cert/server.key")
	if err != nil {
		log.Fatalf("Failed to setup tls: %v", err)
	}

	lis, err := net.Listen("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	acc, err := entity.NewAccount("name", testdata.CPF().String(), testdata.Secret().String())
	if err != nil {
		log.Fatalf("failed to create account: %v", err)
	}

	accRepo := &InMemoryAccountDatabase{Accounts: []entity.Account{acc}}
	trRepo := &InMemoryTransferDatabase{}
	enc := &InMemoryEncrypter{}
	jwt := &InMemoryJWT{}
	accountUsecase := account.NewAccountUsecase(accRepo, enc)
	transferUsecase := transfer.NewTransferUsecase(trRepo, accountUsecase)
	auth := auth.NewAuth(accountUsecase, enc, jwt)

	srv := grpcServer.NewServer(accountUsecase, transferUsecase, auth)

	s := grpc.NewServer(
		grpc.Creds(creds),
		grpc.UnaryInterceptor(grpcServer.AuthInterceptor),
	)

	proto.RegisterBankingServiceServer(s, srv)

	log.Println("Server listening on localhost:9000!")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
