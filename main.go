package main

import (
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/olivere/elastic.v5/config"
	"time"
)

type WsMsg struct{
	Msg string `json:msg`
}
var logMsgs = make(chan *WsMsg,1000)

func StartGetLogServer(logpath string) {
	go ReadLogLoop(logpath)
	WriteLog2Ws()
}

//TODO channel
func ReadLogLoop(logpath string){
	t, _ := tail.TailFile(logpath, tail.Config{Follow: true})
	for line := range t.Lines {
		logMsgs <- &WsMsg{Msg:line.Text}
	}
}

func WriteLog2Ws() {
	for{
		msg := <- logMsgs
		fmt.Println(msg.Msg)
	}
}

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string                `json:"user"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

//const mapping = `
//{
//	"settings":{
//		"number_of_shards": 1,
//		"number_of_replicas": 0
//	},
//	"mappings":{
//		"tweet":{
//			"properties":{
//				"user":{
//					"type":"keyword"
//				},
//				"message":{
//					"type":"text",
//					"store": true,
//					"fielddata": true
//				},
//				"image":{
//					"type":"keyword"
//				},
//				"created":{
//					"type":"date"
//				},
//				"tags":{
//					"type":"keyword"
//				},
//				"location":{
//					"type":"geo_point"
//				},
//				"suggest_field":{
//					"type":"completion"
//				}
//			}
//		}
//	}
//}`


func main(){
	//StartGetLogServer("access.log")

	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	config := &config.Config{URL:"http://127.0.0.1:9200"}
	client, err := elastic.NewClientFromConfig(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	exists, err := client.IndexExists("weberr-2019.10.28").Do(ctx)
	if err != nil{
		fmt.Println(err)
		return
	}

	if !exists{
		createIndex, err := client.CreateIndex("weberr-2019.10.28").Do(ctx)
		if 	err != nil{
			fmt.Println(err)
		}
		if createIndex.Acknowledged{
			//Not acknowledged
		}
	}

	tweet1  := Tweet{User: "olivere", Message: "Take Five", Retweets: 0}
	put1, err := client.Index().
		Index("weberr-2019.10.28").
		Type("tweet").
		Id("1").
		BodyJson(tweet1).
		Do(ctx)

	if err != nil {
		// Handle error
		panic(err)
	}

	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)

	//Get tweet with specified ID
	get1, err := client.Get().
		Index("weberr-2019.10.28").
		Type("tweet").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
		return
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}

	fmt.Println("执行到这里结束")
	return
}
