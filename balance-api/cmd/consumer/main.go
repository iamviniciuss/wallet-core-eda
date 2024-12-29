package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	// httpService "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/http"
	// balance "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/balance"
	// healthcheck "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/web/healthcheck"
	_ "github.com/go-sql-driver/mysql"
	usecase "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/application/use_cases"
	infra "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/worker"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/database"
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

	createTransactionUseCase := usecase.NewCreateTransactionUseCase(uow)

	go infra.QueueRunner(infra.QueueRunnerInput{
		CreateTransactionUseCase: createTransactionUseCase,
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	http.HandleFunc("/balance/", func(w http.ResponseWriter, r *http.Request) {
		// accountId := "54964f6c-01f3-4207-85c6-4722022adf95"
		path := r.URL.Path

		parts := strings.Split(path, "/")

		if len(parts) < 3 {
			http.Error(w, "URL inválida", http.StatusBadRequest)
			return
		}

		accountId := parts[2]

		balanceValue, err := usecase.NewGetBalanceByIdUseCase(uow).Execute(context.Background(), accountId)

		if err != nil {
			w.WriteHeader(http.StatusPreconditionFailed)
			fmt.Fprintf(w, "OK")
			return
		}

		balanceValueAsString := fmt.Sprintf("%f", balanceValue)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, balanceValueAsString)
	})

	http.ListenAndServe(":3003", nil)
}
