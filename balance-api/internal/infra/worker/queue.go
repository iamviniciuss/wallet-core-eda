package worker

import (
	usecase "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/application/use_cases"
	queue "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/infra/broker"
)

type QueueRunnerInput struct {
	CreateTransactionUseCase *usecase.CreateTransactionUseCase
}

func QueueRunner(input QueueRunnerInput) {
	queue.NewSQSMessageBroker(input.CreateTransactionUseCase)
}
