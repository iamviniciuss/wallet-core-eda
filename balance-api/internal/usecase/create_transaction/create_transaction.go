package create_transaction

import (
	"context"
	"fmt"

	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/entity"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/gateway"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
)

type CreateTransactionInputDTO struct {
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type CreateTransactionOutputDTO struct {
	ID            string  `json:"id"`
	AccountIDFrom string  `json:"account_id_from"`
	AccountIDTo   string  `json:"account_id_to"`
	Amount        float64 `json:"amount"`
}

type BalanceUpdatedOutputDTO struct {
	AccountIDFrom        string  `json:"account_id_from"`
	AccountIDTo          string  `json:"account_id_to"`
	BalanceAccountIDFrom float64 `json:"balance_account_id_from"`
	BalanceAccountIDTo   float64 `json:"balance_account_id_to"`
}

type CreateTransactionUseCase struct {
	Uow uow.UowInterface
}

func NewCreateTransactionUseCase(
	Uow uow.UowInterface,
) *CreateTransactionUseCase {
	return &CreateTransactionUseCase{
		Uow: Uow,
	}
}

func (uc *CreateTransactionUseCase) Execute(ctx context.Context, input BalanceUpdatedOutputDTO) (*CreateTransactionOutputDTO, error) {
	output := &CreateTransactionOutputDTO{}

	err := uc.Uow.Do(ctx, func(_ *uow.Uow) error {

		accountRepository := uc.getAccountRepository(ctx)

		accountFrom := &entity.Account{
			ID:      input.AccountIDFrom,
			Balance: input.BalanceAccountIDFrom,
		}

		accountTo := &entity.Account{
			ID:      input.AccountIDTo,
			Balance: input.BalanceAccountIDTo,
		}

		_, err := uc.UpsertAccount(ctx, accountFrom)
		if err != nil {
			return err
		}

		_, err = uc.UpsertAccount(ctx, accountTo)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}

		err = accountRepository.UpdateBalance(accountTo)
		fmt.Println("Chegou ao final")
		return err
	})

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return output, nil
}

func (uc *CreateTransactionUseCase) UpsertAccount(ctx context.Context, account *entity.Account) (*entity.Account, error) {

	accountRepository := uc.getAccountRepository(ctx)

	accountFrom, err := accountRepository.FindByID(account.ID)
	if accountFrom == nil || err != nil {
		errSave := accountRepository.Save(account)

		if errSave != nil {
			return nil, errSave
		}
	}

	accountFrom, err = accountRepository.FindByID(account.ID)

	return accountFrom, err
}

func (uc *CreateTransactionUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}
