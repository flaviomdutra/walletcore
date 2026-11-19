package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/flaviomdutra/walletcore/internal/database"
	"github.com/flaviomdutra/walletcore/internal/event"
	"github.com/flaviomdutra/walletcore/internal/event/handler"
	"github.com/flaviomdutra/walletcore/internal/usecase/create_account"
	"github.com/flaviomdutra/walletcore/internal/usecase/create_client"
	"github.com/flaviomdutra/walletcore/internal/usecase/create_transaction"
	"github.com/flaviomdutra/walletcore/internal/web"
	"github.com/flaviomdutra/walletcore/internal/web/webserver"
	"github.com/flaviomdutra/walletcore/pkg/events"
	"github.com/flaviomdutra/walletcore/pkg/kafka"
	"github.com/flaviomdutra/walletcore/pkg/uow"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbHost := envOrDefault("DB_HOST", "mysql")
	dbPort := envOrDefault("DB_PORT", "3306")
	dbUser := envOrDefault("DB_USER", "root")
	dbPassword := envOrDefault("DB_PASSWORD", "root")
	dbName := envOrDefault("DB_NAME", "wallet")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers": "kafka:29092",
		"group.id":          "wallet",
	}

	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()

	clientDB := database.NewClientDB(db)
	accountDB := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)
	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})
	createClientUseCase := create_client.NewCreateClientUseCase(clientDB)
	createAccountUseCase := create_account.NewCreateAccountUseCase(accountDB, clientDB)
	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(
		uow,
		eventDispatcher,
		transactionCreatedEvent,
	)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transaction", transactionHandler.CreateTransaction)

	webserver.Start()
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
