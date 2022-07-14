package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/thalissonfelipe/banking/banking/tests/testdata"
	proto "github.com/thalissonfelipe/banking/gen/banking/v1"
)

func main() {
	rand.Seed(time.Now().Unix())

	creds, err := credentials.NewClientTLSFromFile("../cert/server.crt", "localhost")
	if err != nil {
		log.Fatalf("cert load error: %v", err)
	}

	conn, err := grpc.Dial(":9000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Panicf("did not connect: %v", err)
	}
	defer conn.Close()

	client := proto.NewBankingServiceClient(conn)

	accounts, err := client.ListAccounts(context.Background(), &proto.ListAccountsRequest{})
	if err != nil {
		log.Panicf("Error when calling ListAccounts: %v", err)
	}

	log.Printf("response from server: %s", accounts.Accounts)

	balance, err := client.GetAccountBalance(context.Background(), &proto.GetAccountBalanceRequest{
		AccountId: accounts.Accounts[0].Id,
	})
	if err != nil {
		log.Panicf("Error when calling GetAccountBalance: %v", err)
	}

	log.Printf("response from server: %d", balance.Balance)

	id, err := client.CreateAccount(context.Background(), &proto.CreateAccountRequest{
		Name:   "Sousa",
		Cpf:    testdata.CPF().String(),
		Secret: testdata.Secret().String(),
	})
	if err != nil {
		log.Panicf("Error when calling CreateAccount: %v", err)
	}

	log.Printf("response from server: %s", id.Id)

	accounts, err = client.ListAccounts(context.Background(), &proto.ListAccountsRequest{})
	if err != nil {
		log.Panicf("Error when calling ListAccounts: %v", err)
	}

	log.Printf("response from server: %s", accounts.Accounts)

	token, err := client.Login(context.Background(), &proto.LoginRequest{
		Cpf:    accounts.Accounts[0].Cpf,
		Secret: testdata.CPF().String(),
	})
	if err != nil {
		log.Panicf("Error when calling Login: %v", err)
	}

	log.Printf("response from server: %s", token.Token)

	md := metadata.Pairs("authorization", token.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	transfers, err := client.ListTransfers(ctx, &proto.ListTransfersRequest{})
	if err != nil {
		log.Panicf("Error when calling ListTransfers: %v", err)
	}

	log.Printf("response from server: %s", transfers.Transfers)

	amount := 100

	_, err = client.PerformTransfer(ctx, &proto.PerformTransferRequest{
		AccountDestinationId: id.Id,
		Amount:               int64(amount),
	})
	if err != nil {
		log.Panicf("Error when calling PerformTransfer: %v", err)
	}

	transfers, err = client.ListTransfers(ctx, &proto.ListTransfersRequest{})
	if err != nil {
		log.Panicf("Error when calling ListTransfers: %v", err)
	}

	log.Printf("response from server: %s", transfers.Transfers)
}
