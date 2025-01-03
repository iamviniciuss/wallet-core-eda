package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	usecase "github.com/iamviniciuss/wallet-core-eda/balance-api/internal/application/use_cases"
)

type KafkaAdapter struct {
}

type TransactionMessage struct {
	Name    string                          `json:"Name"`
	Payload usecase.BalanceUpdatedOutputDTO `json:"Payload"`
}

type ConsumeTopic struct {
	TopicName string
	Handler   func(message ckafka.Message) error
}

func NewSQSMessageBroker(createTransactionUseCase *usecase.CreateTransactionUseCase) *KafkaAdapter {
	KAFKA_URL := os.Getenv("KAFKA_URL")
	fmt.Println("KAFKA_URL", KAFKA_URL)
	configMap := ckafka.ConfigMap{
		"bootstrap.servers":   KAFKA_URL,
		"group.id":            "wallet",
		"auto.offset.reset":   "earliest",
		"api.version.request": false,
	}

	topics := []ConsumeTopic{
		{
			TopicName: "balances",
			Handler: func(message ckafka.Message) error {
				fmt.Println("Mensagem recebida: %v\n", string(message.Value))
				fmt.Println("\n---------------------------\n")

				var dto TransactionMessage

				err := json.Unmarshal(message.Value, &dto)

				if err != nil {
					return err
				}

				fmt.Println(dto.Name)
				fmt.Println(dto.Payload.AccountIDFrom)
				fmt.Println(dto.Payload.AccountIDTo)
				fmt.Println(dto.Payload.BalanceAccountIDFrom)
				fmt.Println(dto.Payload.BalanceAccountIDTo)
				ctx := context.Background()

				_, err = createTransactionUseCase.Execute(ctx, usecase.BalanceUpdatedOutputDTO{
					AccountIDFrom:        dto.Payload.AccountIDFrom,
					BalanceAccountIDFrom: dto.Payload.BalanceAccountIDFrom,
					AccountIDTo:          dto.Payload.AccountIDTo,
					BalanceAccountIDTo:   dto.Payload.BalanceAccountIDTo,
				})

				fmt.Println(err)

				return err
			},
		},
	}
	consumer, err := ckafka.NewConsumer(&configMap)

	if err != nil {
		fmt.Printf("Erro ao criar o consumidor: %v\n", err)
		panic(err)
	}

	defer consumer.Close()

	for _, topicConfig := range topics {
		err := consumer.SubscribeTopics([]string{topicConfig.TopicName}, nil)
		if err != nil {
			panic(err)
		}
		go StartConsume(topicConfig, consumer)
	}

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan

	fmt.Println("Saindo do consumidor Kafka")

	return &KafkaAdapter{}
}

func (svc *KafkaAdapter) Consumer(queue string, consumer ConsumerHandlers) bool {
	return false
}

func StartConsume(topicConfig ConsumeTopic, consumer *ckafka.Consumer) {

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Printf("Consumindo mensagens do tópico: %s\n", topicConfig.TopicName)

	run := true
	for run {
		select {
		case sig := <-sigchan:
			fmt.Println("Sinal recebido: %v\n", sig)
			// signal.Stop(sigchan)

			panic("stop consumer")
			//
		default:
			ev := consumer.Poll(100) // Poll timeout de 100 ms
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *ckafka.Message:

				err := topicConfig.Handler(*e)
				if err != nil {
					fmt.Println("ErroHandler: %v\n", err.Error())
				}
			case ckafka.Error:
				fmt.Println("Erro: %v\n", e.Error())
				panic(e)
			}
		}
	}
}
