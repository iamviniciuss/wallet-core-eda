package get_balance_by_account_id

import (
	"context"

	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/entity"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/gateway"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/pkg/uow"
)

type GetBalanceByIdUseCase struct {
	Uow uow.UowInterface
}

func NewGetBalanceByIdUseCase(
	Uow uow.UowInterface,
) *GetBalanceByIdUseCase {
	return &GetBalanceByIdUseCase{
		Uow: Uow,
	}
}

func (uc *GetBalanceByIdUseCase) Execute(ctx context.Context, accountID string) (*entity.Account, error) {
	accountRepository := uc.getAccountRepository(ctx)

	account, err := accountRepository.FindByID(accountID)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (uc *GetBalanceByIdUseCase) getAccountRepository(ctx context.Context) gateway.AccountGateway {
	repo, err := uc.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}
