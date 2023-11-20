package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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

	mysql_migrations_connection, err := repository.NewConnection()

	if err != nil {
		panic(err)
	}

	CreateDatabaseIfNotExists(mysql_migrations_connection)
	defer mysql_migrations_connection.Close()

	migrationRepo := repository.NewMigrationRepositoryMySQL(mysql_migrations_connection)
	migrationRepo.CreateCollectionIfNotExists("migrations")

	mysql_connection, err := NewConnection()
	if err != nil {
		panic(err)
	}
	defer mysql_connection.Close()

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
	CREATE TABLE wallet.clients (
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
	CREATE TABLE wallet.accounts (
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
	CREATE TABLE wallet.transactions (
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
	createTableQuery := fmt.Sprintf(`CREATE SCHEMA wallet DEFAULT CHARACTER SET utf8;`)
	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		fmt.Println("not created the WALLET database", err.Error())
		return err
	}

	fmt.Println("created the WALLET database")

	_, err2 := mysqlDb.Exec("USE wallet;")
	if err2 != nil {
		fmt.Println("error connecting to 'wallet' database:", err2.Error())
		return err2
	}

	return nil
}

var (
	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	hostname = os.Getenv("MYSQL_HOST")
	port     = os.Getenv("MYSQL_PORT")
	dbName   = os.Getenv("MYSQL_DATABASE")
)

func NewConnection() (*sql.DB, error) {
	dbPort, err := strconv.Atoi(port)

	if err != nil {
		fmt.Println("Error converting string to an integer:", err)
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
