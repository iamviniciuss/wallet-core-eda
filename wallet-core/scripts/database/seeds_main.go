package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	pkg "github.com/iamviniciuss/golang-migrations/src/pkg"
	"github.com/iamviniciuss/golang-migrations/src/repository"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/database"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/event"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/event/handler"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/pkg/events"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/pkg/kafka"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core/pkg/uow"
	seed "github.com/iamviniciuss/wallet-core-eda/wallet-core/scripts/database/seeds"
)

func main() {
	log.Println("The seeds execution has been started.")

	mysql_connection, err := repository.NewConnection()

	if err != nil {
		panic(err)
	}

	defer mysql_connection.Close()

	migrationRepo := repository.NewMigrationRepositoryMySQL(mysql_connection)
	CreateDatabaseIfNotExists(mysql_connection)
	migrationRepo.CreateCollectionIfNotExists("migrations")
	CreateClientsCollectionIfNotExists(mysql_connection)
	CreateAccountsCollectionIfNotExists(mysql_connection)
	CreateTransactionsCollectionIfNotExists(mysql_connection)

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

func CreateClientsCollectionIfNotExists(mysqlDb *sql.DB) error {

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE clients (
		id varchar(45) NOT NULL,
		name varchar(45) DEFAULT NULL,
		email varchar(45) DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		updated_at datetime DEFAULT NULL,
		PRIMARY KEY (id)
	  );`)

	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateAccountsCollectionIfNotExists(mysqlDb *sql.DB) error {

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE accounts (
		id varchar(45) NOT NULL,
		client_id varchar(45) DEFAULT NULL,
		balance float DEFAULT NULL,
		created_at datetime DEFAULT NULL,
		PRIMARY KEY (id)
	  );`)

	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateTransactionsCollectionIfNotExists(mysqlDb *sql.DB) error {

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE transactions (
		id varchar(45) NOT NULL,
		created_at datetime DEFAULT NULL,
		account_id_from varchar(45) DEFAULT NULL,
		account_id_to varchar(45) DEFAULT NULL,
		amount float DEFAULT NULL,
		PRIMARY KEY (id)
	  );`)

	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		return err
	}

	return nil
}

func CreateDatabaseIfNotExists(mysqlDb *sql.DB) error {

	createTableQuery := fmt.Sprintf(`CREATE DATABASE IF NOT EXISTS wallet;`)
	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		fmt.Println("Não criou o banco de dados", err.Error())
		return err
	}

	return nil
}
