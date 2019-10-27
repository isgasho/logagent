package main

import (
	"context"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/olivere/elastic.v5/config"
	"time"
)

// ErrMonitor is a structure used for serializing/deserializing data in Elasticsearch.
type ErrMonitor struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"tweet":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func main(){
	ctx := context.Background()

	config := &config.Config{URL:"http://192.168.85.211:9200"}
	client, err := elastic.NewClientFromConfig(config)
	if err != nil {
		fmt.Println(err)
	}

	exists, err := client.IndexExists("weberr-2019.10.27").Do(ctx)
	if err != nil{
		fmt.Println(err)
	}

	if !exists{
		createIndex, err := client.CreateIndex("weberr-2019.10.27").BodyString(mapping).Do(ctx)
		if 	err != nil{
			fmt.Println(err)
		}
		if createIndex.Acknowledged{
			//Not acknowledged
		}
	}

	errMonitor1 := ErrMonitor{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("weberr-2019.10.27").
		Type("tweet").
		Id("1").
		BodyJson(errMonitor1).
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

}
