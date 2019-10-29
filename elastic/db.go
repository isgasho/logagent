package elastic

import (
	"fmt"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

const (
	host = "http://192.168.85.211:9200"
)

var (
	// elastic 的数据库
	elasticClient *elastic.Client
)

func init() {
	fmt.Println("|foundation|init|db|Init")

	var (
		err error
	)

	cfg := &config.Config{URL: host}
	if elasticClient, err = elastic.NewClientFromConfig(cfg); err != nil {
		fmt.Sprintf("Elastic Client Init Err:%v", err)
		panic(err)
	}
}

func ElasticClient() *elastic.Client {
	return elasticClient
}
