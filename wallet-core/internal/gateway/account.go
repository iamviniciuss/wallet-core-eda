package gateway

import "github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/entity"

type AccountGateway interface {
	Save(account *entity.Account) error
	FindByID(id string) (*entity.Account, error)
	UpdateBalance(account *entity.Account) error
}
