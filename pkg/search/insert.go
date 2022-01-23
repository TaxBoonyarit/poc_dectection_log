package search

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
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

func (e *Elasticsearch) InsertData() {
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
	_, err = e.Elastic.Create("log", id.String(), read)
	if err != nil {
		panic(err)
	}
}
