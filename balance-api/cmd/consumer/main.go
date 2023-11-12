package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	// httpService "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
	// balance "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/balance"
	// healthcheck "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/healthcheck"
	_ "github.com/go-sql-driver/mysql"
	infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/worker"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/database"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/usecase/create_transaction"
)

func main() {
	MYSQL_USER := os.Getenv("MYSQL_USER")
	MYSQL_PASSWORD := os.Getenv("MYSQL_PASSWORD")
	MYSQL_HOST := os.Getenv("MYSQL_HOST")
	MYSQL_PORT := os.Getenv("MYSQL_PORT")
	MYSQL_DATABASE := os.Getenv("MYSQL_DATABASE_BALANCE")

	mysqlhost := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", MYSQL_USER, MYSQL_PASSWORD, MYSQL_HOST, MYSQL_PORT, MYSQL_DATABASE)

	fmt.Println(mysqlhost)

	db, err := sql.Open("mysql", mysqlhost)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	ctx := context.Background()
	uow := uow.NewUow(ctx, db)

	uow.Register("AccountDB", func(tx *sql.Tx) interface{} {
		return database.NewAccountDB(db)
	})

	uow.Register("TransactionDB", func(tx *sql.Tx) interface{} {
		return database.NewTransactionDB(db)
	})

	createTransactionUseCase := create_transaction.NewCreateTransactionUseCase(uow)

	infra.QueueRunner(infra.QueueRunnerInput{
		CreateTransactionUseCase: createTransactionUseCase,
	})

	// http := httpService.NewFiberHttp()
	// healthcheck.HealthCheckRouter(http)
	// balance.BalanceRouter(http, uow)

	// err = http.ListenAndServe(":3003")
	// if err != nil {
	// 	panic(err)
	// }
}