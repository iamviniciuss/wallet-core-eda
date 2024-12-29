package seed

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	gateway "github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/gateway"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/application/use_cases"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/entity"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/internal/infra/database"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/pkg/events"
	"github.com/iamviniciuss/wallet-core-eda/wallet-core-api/pkg/uow"
)

type StartAccounts struct {
	typing             string
	name               string
	version            uint64
	accountGateway     gateway.AccountGateway
	DB                 *sql.DB
	clientGateway      gateway.ClientGateway
	Uow                uow.UowInterface
	eventDispatcher    events.EventDispatcherInterface
	transactionCreated events.EventInterface
	balanceUpdated     events.EventInterface
}

func NewStartAccounts(
	Uow uow.UowInterface,
	eventDispatcher events.EventDispatcherInterface,
	transactionCreated events.EventInterface,
	balanceUpdated events.EventInterface,
	DB *sql.DB,
) *StartAccounts {
	return &StartAccounts{
		version:            1,
		name:               "StartAccounts",
		typing:             "seed",
		Uow:                Uow,
		eventDispatcher:    eventDispatcher,
		transactionCreated: transactionCreated,
		balanceUpdated:     balanceUpdated,
		DB:                 DB,
	}
}

func (onl *StartAccounts) GetName() string {
	return onl.name
}

func (onl *StartAccounts) GetType() string {
	return onl.typing
}

func (onl *StartAccounts) GetVersion() uint64 {
	return onl.version
}

func (onl *StartAccounts) Up() error {
	clientDbRepository := database.NewClientDB(onl.DB)

	client := &entity.Client{
		ID:        "d5a35295-4e15-4a15-99c1-8245b8467a8c",
		Name:      "Vinicius Santos",
		Email:     "vinicius@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Accounts:  []*entity.Account{},
	}

	err := clientDbRepository.Save(client)
	if err != nil {
		return err
	}

	account := &entity.Account{
		ID:        "54964f6c-01f3-4207-85c6-4722022adf95",
		Client:    client,
		ClientID:  client.ID,
		Balance:   1000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = database.NewAccountDB(onl.DB).Save(account)
	if err != nil {
		return err
	}

	// CLIENT 2
	client2 := &entity.Client{
		ID:        "d5a76543-4e15-4a15-99c1-8245b8467v6a",
		Name:      "Camila Santos",
		Email:     "camila@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Accounts:  []*entity.Account{},
	}

	err = clientDbRepository.Save(client2)
	if err != nil {
		return err
	}

	account2 := &entity.Account{
		ID:        "54964f6c-01f3-4207-85c6-8245b8467f2b",
		Client:    client2,
		ClientID:  client2.ID,
		Balance:   2000,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = database.NewAccountDB(onl.DB).Save(account2)
	if err != nil {
		fmt.Println("NewAccountDB2.Save", err.Error())
		return err
	}

	output, err := use_cases.NewCreateTransactionUseCase(onl.Uow, onl.eventDispatcher, onl.transactionCreated, onl.balanceUpdated).Execute(
		context.Background(),
		use_cases.CreateTransactionInputDTO{
			AccountIDFrom: account.ID,
			AccountIDTo:   account2.ID,
			Amount:        100,
		},
	)

	if err != nil {
		return err
	}

	fmt.Println("Transaction ID:", output.ID)
	fmt.Println("Up:", onl.typing, ": ", onl.name, "executed with success")
	return nil
}

func (onl *StartAccounts) Down() error {
	fmt.Println("Down: ", onl.typing, ": ", onl.name, "executed with success")
	return nil
}
