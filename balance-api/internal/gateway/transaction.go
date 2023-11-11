package gateway

import "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
