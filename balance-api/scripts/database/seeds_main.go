package main

import (
	// "context"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/iamviniciuss/golang-migrations/src/repository"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/application/usecase"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/database"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/scripts/database/seeds"
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

	ctx := context.Background()
	uow := uow.NewUow(ctx, mysql_migrations_connection)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(mysql_migrations_connection)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(mysql_migrations_connection)
	})

	mysql_connection, err := NewConnection()
	if err != nil {
		panic(err)
	}
	defer mysql_connection.Close()

	CreateAccountsCollectionIfNotExists(mysql_connection)

	useCase := usecase.NewCreateTransactionUseCase(uow)
	_, err1 := useCase.Execute(context.TODO(), usecase.BalanceUpdatedOutputDTO{
		AccountIDFrom:        "d5a35295-4e15-4a15-99c1-8245b8467a8c",
		AccountIDTo:          "d5a76543-4e15-4a15-99c1-8245b8467v6a",
		BalanceAccountIDFrom: 900,
		BalanceAccountIDTo:   2100,
	})

	if err1 != nil {
		fmt.Println("NewCreateTransactionUseCase", err1.Error())
	}

	seeds.NewStartAccounts(uow, mysql_migrations_connection).Up()

	log.Println("The seeds execution has been successfully completed.")
}

func CreateAccountsCollectionIfNotExists(mysqlDb *sql.DB) error {

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE balance.accounts (
		id varchar(45) NOT NULL,
		balance float DEFAULT NULL
	  );`)

	output, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		fmt.Println("output", output)
		return err
	}

	fmt.Println("output", output)

	return nil
}

func CreateDatabaseIfNotExists(mysqlDb *sql.DB) error {
	createTableQuery := fmt.Sprintf(`CREATE SCHEMA balance DEFAULT CHARACTER SET utf8;`)
	_, err := mysqlDb.Exec(createTableQuery)
	if err != nil {
		log.Println("Did not create the balance database", err.Error())
		return err
	}

	fmt.Println("Created the balance database")

	_, err2 := mysqlDb.Exec("USE balance;")
	if err2 != nil {
		log.Println("Error connecting to 'balance' database:", err2.Error())
		return err2
	}

	return nil
}

var (
	username = os.Getenv("MYSQL_USER")
	password = os.Getenv("MYSQL_PASSWORD")
	hostname = os.Getenv("MYSQL_HOST")
	port     = os.Getenv("MYSQL_PORT")
	dbName   = os.Getenv("MYSQL_DATABASE_BALANCE")
)

func NewConnection() (*sql.DB, error) {
	dbPort, err := strconv.Atoi(port)

	if err != nil {
		fmt.Println("Erro ao converter a string para um número inteiro:", err)
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", username, password, hostname, dbPort, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
