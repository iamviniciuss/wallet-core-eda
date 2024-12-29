package broker

type Message struct {
	Body          string `json:"body"`
	MessageId     string `json:"message_id"`
	ReceiptHandle string `type:"string"`
}

type ConsumerHandlers interface {
	GetTopic() string
	Handler(message *Message) error
	ErrorHandler(err error)
}

type Consumer struct {
	ConsumerHandlers
	queueUrl string
}
