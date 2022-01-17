package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type Filter struct {
	IP        string    `json:"ip"`
	TimeStamp time.Time `json:"timestamp"`
}

type Rule struct {
	Index        string `json:"index"`
	Filter       Filter `json:"filter"`
	Total        int    `json:"total"`
	TimeDuration int    `json:"time_duration"`
}

var rules = []Rule{
	{
		Index: "",
		Filter: Filter{
			IP:        "",
			TimeStamp: time.Now(),
		},
		Total:        0,
		TimeDuration: 0,
	},
}

func main() {
	empJSON, err := json.MarshalIndent(rules, "", " ")
	if err != nil {
		panic(err)
	}

	jsonPretty := fmt.Sprintf(" %s \n", string(empJSON))
	fmt.Printf(" Pretty Json : %s \n", jsonPretty)
}
