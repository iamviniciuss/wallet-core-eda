package worker

import (
	queue "github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/infra/broker"
	"github.com.br/devfullcycle/fc-ms-wallet/balance-api/internal/usecase/create_transaction"
)

type QueueRunnerInput struct {
	CreateTransactionUseCase *create_transaction.CreateTransactionUseCase
}

func QueueRunner(input QueueRunnerInput) {
	queue.NewSQSMessageBroker(input.CreateTransactionUseCase)
}
