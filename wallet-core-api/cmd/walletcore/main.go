package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	_ "github.com/go-sql-driver/mysql"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/event"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/use_cases"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/infra/database"
	handler "github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/infra/event_handler"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/infra/web"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/infra/web/webserver"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/pkg/events"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/pkg/kafka"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/pkg/uow"
)

func main() {
	KAFKA_URL := os.Getenv("KAFKA_URL")
	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE")

	MYSQL_URL_CONN := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE)

	fmt.Println(MYSQL_URL_CONN)

	db, err := sql.Open("mysql", MYSQL_URL_CONN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	configMap := ckafka.ConfigMap{
		"bootstrap.servers":   KAFKA_URL, // Lista de servidores Kafka
		"api.version.request": false,     // Configuração para solicitar a versão da API automaticamente
	}
	kafkaProducer := kafka.NewKafkaProducer(&configMap)

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("TransactionCreated", handler.NewTransactionCreatedKafkaHandler(kafkaProducer))
	eventDispatcher.Register("BalanceUpdated", handler.NewUpdateBalanceKafkaHandler(kafkaProducer))
	transactionCreatedEvent := event.NewTransactionCreated()
	balanceUpdatedEvent := event.NewBalanceUpdated()

	clientDb := database.NewClientDB(db)
	accountDb := database.NewAccountDB(db)

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createTransactionUseCase := use_cases.NewCreateTransactionUseCase(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent)
	createClientUseCase := use_cases.NewCreateClientUseCase(clientDb)
	createAccountUseCase := use_cases.NewCreateAccountUseCase(accountDb, clientDb)

	webserver := webserver.NewWebServer(":8080")

	clientHandler := web.NewWebClientHandler(*createClientUseCase)
	accountHandler := web.NewWebAccountHandler(*createAccountUseCase)
	transactionHandler := web.NewWebTransactionHandler(*createTransactionUseCase)

	webserver.AddHandler("/clients", clientHandler.CreateClient)
	webserver.AddHandler("/accounts", accountHandler.CreateAccount)
	webserver.AddHandler("/transactions", transactionHandler.CreateTransaction)
	webserver.AddHandler("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	fmt.Println("Server is running")
	webserver.Start()

}
