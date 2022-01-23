package main

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/google/uuid"
	"github.com/logzio/logzio-go"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Log struct {
	Timestamp    time.Time `json:"@timestamp"`
	Level        string    `json:"level"`
	LogChain     string    `json:"log_chain"`
	SrcInterface string    `json:"src_interface"`
	Out          string    `json:"out"`
	Proto        string    `json:"proto"`
	SrcIp        string    `json:"src_ip"`
	DstIp        string    `json:"dst_ip"`
	Hostname     string    `json:"hostname"`
	TypeLog      string    `json:"type_log"`
}

func main() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Elasticsearch initializing....")

	l, err := logzio.New(
		"ThLktEByGqDUQvPPaRhzYRUXLqbrfhEV",
		logzio.SetDebug(os.Stderr),
		logzio.SetUrl("https://listener.logz.io:8071"),
		logzio.SetDrainDuration(time.Second*5),
		logzio.SetTempDirectory("myQueue"),
		logzio.SetDrainDiskThreshold(99),
	)
	if err != nil {
		panic(err)
	}

	for {
		body := Log{
			Timestamp:    time.Now(),
			Level:        "info",
			LogChain:     "forward",
			SrcInterface: "ether1",
			Out:          "bridge-server",
			Proto:        "UDP",
			SrcIp:        fmt.Sprintf("192.168.1.%d", rand.Int31n(5)),
			DstIp:        "10.10.100.100",
			Hostname:     "dev",
			TypeLog:      "firewall",
		}

		msg, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		if err = l.Send(msg); err != nil {
			panic(err)
		}

		empJSON, err := json.MarshalIndent(body, "", " ")
		if err != nil {
			panic(err)
		}

		jsonPretty := fmt.Sprintf(" %s \n", string(empJSON))

		var b strings.Builder
		_, err = b.WriteString(jsonPretty)
		if err != nil {
			panic(err)
		}
		read := strings.NewReader(b.String())

		id := uuid.New()
		_, err = client.Create("log", id.String(), read)
		if err != nil {
			panic(err)
		}

		fmt.Println("Insert data success")
		time.Sleep(time.Second * 5)
	}
}
