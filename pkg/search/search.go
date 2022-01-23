package search

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v7"
)

type Elasticsearch struct {
	Elastic *elasticsearch.Client
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
