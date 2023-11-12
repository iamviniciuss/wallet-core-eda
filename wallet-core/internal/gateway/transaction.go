package gateway

import "github.com/iamviniciuss/wallet-core-eda/wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
