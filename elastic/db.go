package elastic

import (
	"fmt"
	"github.com/astaxie/beego"

	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
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

	elasticSniff,_ := beego.AppConfig.Bool("elastic.sniff") //因为在内网测试,索引要关闭探寻器，不然会出错
	url :=beego.AppConfig.DefaultString("elastic.url","http://127.0.0.1:9200")

	cfg := &config.Config{URL: url, Sniff: &elasticSniff}
	if elasticClient, err = elastic.NewClientFromConfig(cfg); err != nil {
		fmt.Sprintf("Elastic Client Init Err:%v", err)
		panic(err)
	}
}

func ElasticClient() *elastic.Client {
	return elasticClient
}
