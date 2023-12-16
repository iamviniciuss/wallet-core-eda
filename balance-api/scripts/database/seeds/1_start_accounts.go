package seeds

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/database"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/entity"
	gateway "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/gateway"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
)

type StartAccounts struct {
	typing         string
	name           string
	version        uint64
	accountGateway gateway.AccountGateway
	DB             *sql.DB
	clientGateway  gateway.ClientGateway
	Uow            uow.UowInterface
}

func NewStartAccounts(
	Uow uow.UowInterface,
	DB *sql.DB,
) *StartAccounts {
	return &StartAccounts{
		version: 1,
		name:    "StartAccounts",
		typing:  "seed",
		Uow:     Uow,
		DB:      DB,
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
	client := &entity.Client{
		ID:        "d5a35295-4e15-4a15-99c1-8245b8467a8c",
		Name:      "Vinicius Santos",
		Email:     "vinicius@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Accounts:  []*entity.Account{},
	}

	account := &entity.Account{
		ID:        "54964f6c-01f3-4207-85c6-4722022adf95",
		Client:    client,
		ClientID:  client.ID,
		Balance:   900,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := database.NewAccountDB(onl.DB).Save(account)
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

	account2 := &entity.Account{
		ID:        "54964f6c-01f3-4207-85c6-8245b8467f2b",
		Client:    client2,
		ClientID:  client2.ID,
		Balance:   2100,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = database.NewAccountDB(onl.DB).Save(account2)
	if err != nil {
		fmt.Println("NewAccountDB2.Save", err.Error())
		return err
	}

	if err != nil {
		return err
	}
	return nil
}

func (onl *StartAccounts) Down() error {
	fmt.Println("Down: ", onl.typing, ": ", onl.name, "executed with success")
	return nil
}
