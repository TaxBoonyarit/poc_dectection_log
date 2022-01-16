package queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

type Producer struct {
	Producer sarama.SyncProducer
}

var (
	brokerList = []string{"localhost:9092"}
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

	fmt.Println("Producer initializing....")

	return &Producer{
		Producer: producer,
	}
}
