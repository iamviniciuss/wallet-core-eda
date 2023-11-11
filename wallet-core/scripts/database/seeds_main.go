package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/internal/database"
	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/internal/event"
	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/internal/event/handler"
	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/pkg/events"
	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/pkg/kafka"
	"github.com.br/devfullcycle/fc-ms-wallet/wallet-core/pkg/uow"
	seed "github.com.br/devfullcycle/fc-ms-wallet/wallet-core/scripts/database/seeds"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	pkg "github.com/iamviniciuss/golang-migrations/src/pkg"
	"github.com/iamviniciuss/golang-migrations/src/repository"
)

func main() {
	log.Println("The seeds execution has been started.")

	mysql_connection, err := repository.NewConnection()

	if err != nil {
		panic(err)
	}

	defer mysql_connection.Close()

	migrationRepo := repository.NewMigrationRepositoryMySQL(mysql_connection)
	migrationRepo.CreateCollectionIfNotExists("migrations")
	migrationRepo.CreateClientsCollectionIfNotExists()
	migrationRepo.CreateAccountsCollectionIfNotExists()
	migrationRepo.CreateTransactionsCollectionIfNotExists()

	KAFKA_URL := os.Getenv("KAFKA_URL")
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

	// clientDb := database.NewClientDB(mysql_connection)
	// accountDb := database.NewAccountDB(mysql_connection)

	ctx := context.Background()
	uow := uow.NewUow(ctx, mysql_connection)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(mysql_connection)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(mysql_connection)
	})

	migrationManager := pkg.NewMigrate(migrationRepo)
	migrationManager.Register(seed.NewStartAccounts(uow, eventDispatcher, transactionCreatedEvent, balanceUpdatedEvent, mysql_connection))

	if err := migrationManager.Up(pkg.AllAvailable); err != nil {
		log.Println("An unhandled error occurred during the seeds execution.")
		panic(err)
	}

	log.Println("The seeds execution has been successfully completed.")
}
