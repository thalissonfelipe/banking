package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	accUsecase "github.com/thalissonfelipe/banking/banking/domain/account/usecase"
	"github.com/thalissonfelipe/banking/banking/domain/entities"
	trUsecase "github.com/thalissonfelipe/banking/banking/domain/transfer/usecase"
	"github.com/thalissonfelipe/banking/banking/domain/vos"
	grpcServer "github.com/thalissonfelipe/banking/banking/gateway/grpc"
	"github.com/thalissonfelipe/banking/banking/services/auth"
	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	proto "github.com/thalissonfelipe/banking/proto/banking"
)

type InMemoryAccountDatabase struct {
	Accounts []entities.Account
}

func (i InMemoryAccountDatabase) GetAccounts(ctx context.Context) ([]entities.Account, error) {
	return i.Accounts, nil
}

func (i InMemoryAccountDatabase) GetBalanceByID(ctx context.Context, id vos.AccountID) (int, error) {
	for _, acc := range i.Accounts {
		if acc.ID == id {
			return acc.Balance, nil
		}
	}

	return 0, entities.ErrAccountDoesNotExist
}

func (i *InMemoryAccountDatabase) CreateAccount(ctx context.Context, account *entities.Account) error {
	for _, acc := range i.Accounts {
		if acc.CPF == account.CPF {
			return entities.ErrAccountAlreadyExists
		}
	}

	i.Accounts = append(i.Accounts, *account)

	return nil
}

func (i InMemoryAccountDatabase) GetAccountByCPF(ctx context.Context, cpf vos.CPF) (entities.Account, error) {
	for _, acc := range i.Accounts {
		if acc.CPF == cpf {
			return acc, nil
		}
	}

	return entities.Account{}, entities.ErrAccountDoesNotExist
}

func (i InMemoryAccountDatabase) GetAccountByID(ctx context.Context, id vos.AccountID) (entities.Account, error) {
	for _, acc := range i.Accounts {
		if acc.ID == id {
			return acc, nil
		}
	}

	return entities.Account{}, entities.ErrAccountDoesNotExist
}

type InMemoryEncrypter struct{}

func (i InMemoryEncrypter) Hash(secret string) ([]byte, error) {
	return []byte(secret), nil
}

func (i InMemoryEncrypter) CompareHashAndSecret(hashedSecret, secret []byte) error {
	return nil
}

type InMemoryTransferDatabase struct {
	Transfers []entities.Transfer
}

func (i InMemoryTransferDatabase) GetTransfers(ctx context.Context, id vos.AccountID) ([]entities.Transfer, error) {
	return i.Transfers, nil
}

func (i *InMemoryTransferDatabase) CreateTransfer(ctx context.Context, transfer *entities.Transfer) error {
	tr := entities.NewTransfer(transfer.AccountOriginID, transfer.AccountDestinationID, transfer.Amount)

	i.Transfers = append(i.Transfers, tr)

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

	acc, err := entities.NewAccount("name", testdata.GetValidCPF().String(), testdata.GetValidSecret().String())
	if err != nil {
		log.Fatalf("failed to create account: %v", err)
	}

	accRepo := &InMemoryAccountDatabase{Accounts: []entities.Account{acc}}
	trRepo := &InMemoryTransferDatabase{}
	enc := &InMemoryEncrypter{}
	accountUsecase := accUsecase.NewAccountUsecase(accRepo, enc)
	transferUsecase := trUsecase.NewTransferUsecase(trRepo, accountUsecase)
	auth := auth.NewAuth(accountUsecase, enc)

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
