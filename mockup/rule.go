package mockup

import "time"

type (
	Rule struct {
		Index     string          `json:"index"`
		Filters   Filter          `json:"filters"`
		Duration  int             `json:"duration"`
		Groups    []string        `json:"groups"`
		Aggregate Aggregate       `json:"aggregate"`
		Size      int             `json:"size"`
		TimeRange StartAndEndDate `json:"time_range"`
	}

	Filter struct {
		IP      string `json:"src_ip"`
		TypeLog string `json:"type_log"`
	}

	Aggregate struct {
		Type      string `json:"type"`
		Condition string `json:"condition"`
		Value     int    `json:"value"`
	}
	StartAndEndDate struct {
		StartDate time.Time `json:"start_end"`
		EndDate   time.Time `json:"end_date"`
	}
)

var Rules = []Rule{
	{
		Index: "log",
		Filters: Filter{
			IP:      "192.168.1.1",
			TypeLog: "firewall",
		},
		Duration: 5,
		Groups:   []string{"src_ip"},
		Aggregate: Aggregate{
			Type:      "count",
			Condition: ">",
			Value:     5,
		},
	},
	//{
	//	Index: "log",
	//	Filters: Filter{
	//		IP:      "192.168.1.2",
	//		TypeLog: "firewall",
	//	},
	//	Duration: 6,
	//	Groups:   []string{"src_ip"},
	//	Aggregate: Aggregate{
	//		Type:      "count",
	//		Condition: ">",
	//		Value:     5,
	//	},
	//},
	//{
	//	Index: "log",
	//	Filters: Filter{
	//		IP:      "192.168.1.3",
	//		TypeLog: "firewall",
	//	},
	//	Duration: 7,
	//	Groups:   []string{"src_ip"},
	//	Aggregate: Aggregate{
	//		Type:      "count",
	//		Condition: ">",
	//		Value:     15,
	//	},
	//},
	//{
	//	Index: "log",
	//	Filters: Filter{
	//		IP:      "192.168.1.4",
	//		TypeLog: "firewall",
	//	},
	//	Duration: 8,
	//	Groups:   []string{"src_ip"},
	//	Aggregate: Aggregate{
	//		Type:      "count",
	//		Condition: "<",
	//		Value:     20,
	//	},
	//},
	//{
	//	Index: "log",
	//	Filters: Filter{
	//		IP:      "192.168.1.5",
	//		TypeLog: "firewall",
	//	},
	//	Duration: 9,
	//	Groups:   []string{"src_ip"},
	//	Aggregate: Aggregate{
	//		Type:      "count",
	//		Condition: ">",
	//		Value:     20,
	//	},
	//},
}
