package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/aclgo/balance/config"
	"github.com/aclgo/balance/infra/delivery/grpc/service"
	"github.com/aclgo/balance/infra/repository"
	"github.com/aclgo/balance/proto"
	"github.com/aclgo/balance/usecase"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {

	cfg := config.NewConfig(".")
	if err := cfg.Load(); err != nil {
		log.Fatal(err)
	}

	clientOptions := options.Client().ApplyURI(cfg.DbUrl)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("mongo.Connect; %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("client.Ping: %v", err)
	}

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Fatalf("client.Disconnect: %v", err)
		}
	}()

	collection := client.Database(cfg.DbName).Collection(cfg.DbCollection)
	collectionTransactions := client.Database(cfg.DbName).Collection("transactions")

	repo := repository.NewMongoRepository(collection, collectionTransactions)
	if err := repo.EnsureIndexes(ctx); err != nil {
		log.Fatalf("repo.EnsureIndexes: %w", err)
	}

	mu := sync.Mutex{}

	createUC := usecase.NewWalletCreateUC(repo)
	creditUC := usecase.NewWalletCreditUC(repo, &mu)
	debitUC := usecase.NewWalletDebitUC(repo, &mu)
	getByAccountUC := usecase.NewWalletGetByAccountUC(repo)
	createTransaction := usecase.NewRegisterTransactionUC(repo)

	svcGRPC := service.NewGrpcService(createUC, creditUC, debitUC, getByAccountUC, createTransaction)

	listen, err := net.Listen("tcp", ":50056")

	server := grpc.NewServer()

	proto.RegisterWalletServiceServer(server, svcGRPC)

	log.Printf("server start port 50056\n")
	if err := server.Serve(listen); err != nil {
		log.Fatalf("server.Serve: %v", err)
	}
}
