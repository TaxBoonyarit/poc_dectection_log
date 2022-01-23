package main

import (
	"context"
	"elasticsearch/mockup"
	"elasticsearch/pkg/notification"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/tidwall/gjson"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type (
	Elasticsearch struct {
		Elastic *elasticsearch.Client
	}

	Buckets struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	}
)

var (
	elastic           *elasticsearch.Client
	topic             = "POC_LOG"
	brokerList        = []string{"localhost:9092"}
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := brokerList
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}

	// init elasticsearch
	elasticService := NewElasticsearch()
	elastic = elasticService.Elastic

	fmt.Println("Consumer initializing....")

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	doneCh := make(chan struct{})

	defer func() {
		if err = consumer.Close(); err != nil {
			panic(err)
		}
	}()

	go func() {
		for {
			select {
			case err = <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				*messageCountStart++
				log.Println("Received messages", string(msg.Value))
				var rule mockup.Rule
				_ = json.Unmarshal(msg.Value, &rule)
				SendData(rule)

			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")
}

func NewElasticsearch() *Elasticsearch {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Elasticsearch initializing....")

	return &Elasticsearch{
		Elastic: client,
	}
}

func SendData(rule mockup.Rule) {
	var mapResp map[string]interface{}
	ctx := context.Background()
	duration := time.Duration(rule.Duration) * time.Minute
	startDate := time.Now().Add(-duration).Format(time.RFC3339)
	endDate := time.Now().Format(time.RFC3339)
	query := fmt.Sprintf(`{
			  "size" : %d,
			  "aggregations": {
				"count": {
				  "terms": {
					"field": "src_ip.keyword"
				  }
				}
			  },
			  "query": {
				"bool": {
				  "must": [
					{
					  "match": {
						"type_log": "%s"
					  }
					},
					{
					  "match": {
						"src_ip": "%s"
					  }
					},
					{
					  "range": {
						"@timestamp": {
						  "gte": "%s",
						  "lte": "%s"
						}
					  }
					}
				  ]
				}
			  }
			}`, rule.Size, rule.Filters.TypeLog, rule.Filters.IP, startDate, endDate)

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

	fmt.Println(query)

	if res.StatusCode == http.StatusOK {
		var checkData bool
		if err = json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		result, err := json.Marshal(&mapResp)
		if err != nil {
			panic(err)
		}

		buckets := []Buckets{}
		aggregate := gjson.Get(string(result), "aggregations.count.buckets")

		if err = json.Unmarshal([]byte(aggregate.String()), &buckets); err != nil {
			panic(err)
		}

		for _, v := range buckets {
			if conditionRule(rule.Aggregate.Condition, v.DocCount, rule.Aggregate.Value) {
				checkData = true
			}
		}

		if checkData {
			fmt.Println("found data : ", rule.Filters.IP)
			msg := fmt.Sprintf("\nDetection found IP : %s \nType Log : %s \n\n Condition : \n Time Duration : %d minute \n Total : %d \n ", rule.Filters.IP, rule.Filters.TypeLog, rule.Duration, buckets[0].DocCount)
			msg += fmt.Sprintf("Start Date : %v \n End Date : %v \n", startDate, endDate)
			msg += fmt.Sprintf("Condtion : %d %s %d", buckets[0].DocCount, rule.Aggregate.Condition, rule.Aggregate.Value)
			notification.SendLine(msg)
		}
	}

	defer res.Body.Close()
}

func conditionRule(operator string, n1, n2 int) bool {
	switch operator {
	case ">":
		return n1 > n2
	case "<":
		return n1 < n2
	default:
		return false
	}
}
