package main

import (
	"elasticsearch/mockup"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/robfig/cron"
	"log"
	"time"
)

type (
	Producer struct {
		Producer sarama.SyncProducer
	}
)

var (
	brokerList = []string{"localhost:9092"}
	producer   sarama.SyncProducer
	minute     = 1
	topic      = "POC_LOG"
)

func main() {
	// init producer
	p := NewProducer()
	producer = p.Producer

	// init schedule
	c := cron.New()
	_ = c.AddFunc("@every 1m", checkRule)
	c.Start()

	for {
		time.Sleep(time.Second * 10)
	}
}

func checkRule() {
	if minute == 60 {
		minute = 1
	}

	for _, rule := range mockup.Rules {
		if minute%rule.Duration == 0 {
			out, err := json.Marshal(rule)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Send message : %v \n", rule)
			msg := &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.StringEncoder(out),
			}

			_, _, err = producer.SendMessage(msg)
			if err != nil {
				log.Panic(err)
			}
		}
	}
	minute++
}

func NewProducer() *Producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%v ,Producer initializing....", time.Now().Format(time.RFC3339))

	return &Producer{
		Producer: producer,
	}
}
