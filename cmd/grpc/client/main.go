package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/thalissonfelipe/banking/pkg/tests/testdata"
	proto "github.com/thalissonfelipe/banking/proto/banking"
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

	accounts, err := client.GetAccounts(context.Background(), &proto.ListAccountsRequest{})
	if err != nil {
		log.Panicf("Error when calling GetAccounts: %v", err)
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
		Cpf:    testdata.GetValidCPF().String(),
		Secret: testdata.GetValidSecret().String(),
	})
	if err != nil {
		log.Panicf("Error when calling CreateAccount: %v", err)
	}

	log.Printf("response from server: %s", id.Id)

	accounts, err = client.GetAccounts(context.Background(), &proto.ListAccountsRequest{})
	if err != nil {
		log.Panicf("Error when calling GetAccounts: %v", err)
	}

	log.Printf("response from server: %s", accounts.Accounts)

	token, err := client.Login(context.Background(), &proto.LoginRequest{
		Cpf:    accounts.Accounts[0].Cpf,
		Secret: testdata.GetValidCPF().String(),
	})
	if err != nil {
		log.Panicf("Error when calling Login: %v", err)
	}

	log.Printf("response from server: %s", token.Token)

	md := metadata.Pairs("authorization", token.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	transfers, err := client.GetTransfers(ctx, &proto.ListTransfersRequest{})
	if err != nil {
		log.Panicf("Error when calling GetTransfers: %v", err)
	}

	log.Printf("response from server: %s", transfers.Transfers)

	amount := 100

	_, err = client.CreateTransfer(ctx, &proto.CreateTransferRequest{
		AccountDestinationId: id.Id,
		Amount:               int64(amount),
	})
	if err != nil {
		log.Panicf("Error when calling CreateTransfer: %v", err)
	}

	transfers, err = client.GetTransfers(ctx, &proto.ListTransfersRequest{})
	if err != nil {
		log.Panicf("Error when calling GetTransfers: %v", err)
	}

	log.Printf("response from server: %s", transfers.Transfers)
}
