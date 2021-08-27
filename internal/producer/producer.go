package producer

import (
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Send(message EventMessage) error
}

type producer struct {
	bus   sarama.SyncProducer
	topic string
}

func (p *producer) Send(msg EventMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Err(err).Msg("failed encode to json:")
		return err
	}

	encodedMessage := p.prepareMessage(bytes)
	_, _, err = p.bus.SendMessage(encodedMessage)
	return err
}

func (p *producer) prepareMessage(bytes []byte) *sarama.ProducerMessage {
	encodedMessage := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(p.topic),
		Value:     sarama.StringEncoder(bytes),
		Partition: -1,
	}
	return encodedMessage
}

func NewProducer(brokers []string, topic string) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}
	prod := producer{bus: syncProducer, topic: topic}
	return &prod, nil
}
