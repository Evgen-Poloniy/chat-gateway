package producer

// MessengerProducer represents interface for work with message broker
type MessengerProducer interface {
}

type Producer struct {
	MessengerProducer
}

func NewProducer() *Producer {
	return &Producer{}
}
