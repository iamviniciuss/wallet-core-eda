package infra

type QueueAdapter interface {
	Consumer(queue string, consumer ConsumerHandlers) bool
}
