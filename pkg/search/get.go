package search

import (
	"context"
	"elasticsearch/mockup"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var mapResp map[string]interface{}

func (e *Elasticsearch) GetData(rule mockup.Rule) {
	ctx := context.Background()

	//	var query = fmt.Sprintf(`{
	//   "query" : {
	//       "term" : {
	//           "name" : {
	//               "value" : "%s"
	//           }
	//       }
	//   },
	//	"size" : %d
	//}`, rule.Filter, rule.Total)

	query := "test"
	var b strings.Builder

	_, err := b.WriteString(query)
	if err != nil {
		panic(err)
	}

	read := strings.NewReader(b.String())
	res, err := e.Elastic.Search(
		e.Elastic.Search.WithContext(ctx),
		e.Elastic.Search.WithIndex(rule.Index),
		e.Elastic.Search.WithBody(read),
		e.Elastic.Search.WithTrackTotalHits(true),
		e.Elastic.Search.WithPretty(),
	)

	if err != nil {
		panic(err)
	}

	if res.StatusCode == http.StatusOK {
		fmt.Println(res.StatusCode)
		var checkData bool
		if err = json.NewDecoder(res.Body).Decode(&mapResp); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		}

		for _, hit := range mapResp["hits"].(map[string]interface{})["hits"].([]interface{}) {
			_ = hit.(map[string]interface{})
			checkData = true
		}

		if checkData {
			// alert data
		}
	}

	defer res.Body.Close()
}
