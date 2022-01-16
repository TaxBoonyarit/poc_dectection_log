package queue

import (
	"github.com/Shopify/sarama"
	"log"
)

type Consumer struct {
	Consumer sarama.Consumer
}

func NewConsumer() Consumer {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := brokerList

	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}

	return Consumer{
		Consumer: master,
	}
}
