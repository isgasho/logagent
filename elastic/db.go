package elastic

import (
	"fmt"

	"github.com/astaxie/beego/logs"

	"gopkg.in/olivere/elastic.v5"
	"gopkg.in/olivere/elastic.v5/config"
)

const (
	host = "http://127.0.0.1:9200"
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

	cfg := &config.Config{}
	if elasticClient, err = elastic.NewClientFromConfig(cfg); err != nil {
		logs.Error("Elastic Client Init Err:%v", err)
		panic(err)
	}
}

func ElasticClient() *elastic.Client {
	return elasticClient
}
