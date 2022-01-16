package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/juunini/simple-go-line-notify/notify"
	"github.com/robfig/cron"
	"gopkg.in/alecthomas/kingpin.v2"
	"log"
	"net/http"
	"net/smtp"
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
	e                 *elasticsearch.Client
	minute            = 1
	ctx               = context.Background()
	brokerList        = []string{"localhost:9092"}
	topic             = "testLogs02"
	messageCountStart = kingpin.Flag("messageCountStart", "Message counter start from:").Int()
)

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println(err)
	}

	e = client

	// consumer
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	brokers := brokerList
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err := master.Close(); err != nil {
			log.Panic(err)
		}
	}()

	consumer, err := master.ConsumePartition(topic, 0, sarama.OffsetOldest)
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
			case err := <-consumer.Errors():
				log.Println(err)
			case msg := <-consumer.Messages():
				// have messages
				*messageCountStart++
				log.Println("Received messages", string(msg.Key), string(msg.Value))

				//covert to struct
				var rule Rule
				_ = json.Unmarshal(msg.Value, &rule)

				getDataFoElasticsSearch(rule)
			case <-signals:
				log.Println("Interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	log.Println("Processed", *messageCountStart, "messages")

	//for {
	//	time.Sleep(time.Second * 10)
	//}
}

func schedules() {
	c := cron.New()
	_ = c.AddFunc("0 * * * *", checkRule)
	c.Start()
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

	//data := make(chan interface{}, 1)
	//
	//go func() {
	//	data <- rule
	//	close(data)
	//}()

	//	for res := range data {
	//		covertData := res.(Rule)
	//		var query = fmt.Sprintf(`{
	//    "query" : {
	//        "term" : {
	//            "name" : {
	//                "value" : "%s"
	//            }
	//        }
	//    },
	//	"size" : %d
	//}`, covertData.Filter, covertData.Total)
	//		getDataFoElasticsSearch(query, covertData)
	//	}
}

func getDataFoElasticsSearch(rule Rule) {
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

	res, err := e.Search(
		e.Search.WithContext(ctx),
		e.Search.WithIndex(rule.Index),
		e.Search.WithBody(read),
		e.Search.WithTrackTotalHits(true),
		e.Search.WithPretty(),
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

func sendEmail() {
	from := "inwtax@gmail.com"
	password := "0554861712"
	to := []string{"boonyarit.b@securitypitch.com"}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := []byte("Hello Test.")

	// Create authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send actual message
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}
