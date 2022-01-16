package queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
)

type Producer struct {
	Producer sarama.SyncProducer
}

var (
	brokerList        = []string{"localhost:9092"}
	topic             = "test"
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func NewProducer() *Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Panic(err)
		}
	}()
	fmt.Println("Producer initializing....")

	return &Producer{
		Producer: producer,
	}
}

func (p *Producer) SendData(msg string) {
	m := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	_, _, err := p.Producer.SendMessage(m)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Send message success")
}
