package main

import (
	"encoding/json"
	"fmt"
	"github.com/logzio/logzio-go"
	"math/rand"
	"os"
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
			SrcIp:        fmt.Sprintf("192.168.1.%d", rand.Int31n(10)),
			DstIp:        "10.10.100.100",
			Hostname:     "dev",
			TypeLog:      "firewall",
		}

		msg, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}

		err = l.Send(msg)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Second * 1)
	}

}
