package broker

type Broker struct {
	queue QueueAdapter
}

func NewBroker(adapter QueueAdapter) *Broker {
	return &Broker{
		queue: adapter,
	}
}

func (b *Broker) InitConsumers(consumers []Consumer) {
	for _, consumer := range consumers {
		go b.queue.Consumer(consumer.GetTopic(), consumer)
	}
}
