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
	sniff := false //因为在内网测试,索引要关闭探寻器，不然会出错
	cfg := &config.Config{URL: host, Sniff: &sniff}
	if elasticClient, err = elastic.NewClientFromConfig(cfg); err != nil {
		fmt.Sprintf("Elastic Client Init Err:%v", err)
		panic(err)
	}
}

func ElasticClient() *elastic.Client {
	return elasticClient
}
