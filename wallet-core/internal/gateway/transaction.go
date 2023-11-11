package gateway

import "github.com.br/devfullcycle/fc-ms-wallet/wallet-core/internal/entity"

type TransactionGateway interface {
	Create(transaction *entity.Transaction) error
}
