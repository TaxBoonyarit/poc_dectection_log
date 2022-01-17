package main

import (
	"context"
	"elasticsearch/pkg/queue"
	"elasticsearch/pkg/search"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/juunini/simple-go-line-notify/notify"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

type (
	Rule struct {
		Index        string
		Filter       string
		Total        int
		TimeDuration int
	}
)

var (
	mapResp           map[string]interface{}
	elastic           *elasticsearch.Client
	producer          sarama.SyncProducer
	minute            = 1
	ctx               = context.Background()
	topic             = "test666"
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	// init elasticsearch
	setupElasticsearch := search.NewElasticsearch()
	elastic = setupElasticsearch.Elastic

	// init producer
	p := queue.NewProducer()
	producer = p.Producer

	// init consumer
	c := queue.NewConsumer()
	consumer, err := c.Consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	go func() {
		schedules()
		for {
			select {
			case err = <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				// have messages
				*messageCountStart++
				log.Println("Received messages", string(msg.Key), string(msg.Value))

				//covert to struct
				var rule Rule
				_ = json.Unmarshal(msg.Value, &rule)
				getDataForElasticsSearch(rule)
			case <-signals:
				consumer.Close()
				producer.Close()
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")
}

func schedules() {
	c := cron.New()
	_ = c.AddFunc("0 * * * *", checkRule)
	c.Start()
	fmt.Println("Schedule initializing....")
}

func checkRule() {
	if minute == 60 {
		minute = 1
	}

	rules := []Rule{}
	rules = append(rules, Rule{
		Index:        "products",
		Filter:       "lobster",
		Total:        1,
		TimeDuration: 1,
	})

	rules = append(rules, Rule{
		Index:        "products",
		Filter:       "Beets",
		Total:        1,
		TimeDuration: 1,
	})

	rules = append(rules, Rule{
		Index:        "products",
		Filter:       "testt2",
		Total:        10,
		TimeDuration: 1,
	})

	rules = append(rules, Rule{
		Index:        "products",
		Filter:       "wine",
		Total:        10,
		TimeDuration: 5,
	})
	for _, rule := range rules {
		if minute%rule.TimeDuration == 0 {
			sendDataToKafka(rule)
		}
	}
	minute++
}

func sendDataToKafka(rule Rule) {
	out, err := json.Marshal(rule)
	if err != nil {
		panic(err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(out),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		log.Panic(err)
	}
}

func getDataForElasticsSearch(rule Rule) {
	var query = fmt.Sprintf(`{
   "query" : {
       "term" : {
           "name" : {
               "value" : "%s"
           }
       }
   },
	"size" : %d
}`, rule.Filter, rule.Total)

	var b strings.Builder

	_, err := b.WriteString(query)
	if err != nil {
		panic(err)
	}

	read := strings.NewReader(b.String())
	res, err := elastic.Search(
		elastic.Search.WithContext(ctx),
		elastic.Search.WithIndex(rule.Index),
		elastic.Search.WithBody(read),
		elastic.Search.WithTrackTotalHits(true),
		elastic.Search.WithPretty(),
	)

	if err != nil {
		panic(err)
	}

	if res.StatusCode == http.StatusOK {
		var checkData bool
		if err = json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {
			_ = hit.(map[string]interface{})
			checkData = true
		}

		if checkData {
			sendLine(rule)
		}
	}

	defer res.Body.Close()
}

func sendLine(rule Rule) {
	accessToken := "efpb0Pc6Pyzhcn8745jFQ3zWCNftjfI2u4UsKeLJX3m"
	message := fmt.Sprintf("\n found data index : %s \n filter name : %s \n  Time duration : %d minute", rule.Index, rule.Filter, rule.TimeDuration)
	if err := notify.SendText(accessToken, message); err != nil {
		panic(err)
	}
}
