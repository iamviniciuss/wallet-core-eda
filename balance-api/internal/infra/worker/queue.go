package worker

import (
	queue "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/broker"
	"github.com/iamviniciuss/wallet-core-eda/balance-api/internal/usecase/create_transaction"
)

type QueueRunnerInput struct {
	CreateTransactionUseCase *create_transaction.CreateTransactionUseCase
}

func QueueRunner(input QueueRunnerInput) {
	queue.NewSQSMessageBroker(input.CreateTransactionUseCase)
}
