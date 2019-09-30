package main

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

func main() {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		// Handle error
		panic(err)
	}
	info, code, err := client.Ping("http://127.0.0.1:9200").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	resp, err := client.IndexExists(".kibana_").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Println(resp)
}
