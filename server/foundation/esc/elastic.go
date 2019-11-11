package esc

import (
	"context"
	"fmt"

	"github.com/astaxie/beego/logs"
	"github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

var elasticSourceMap = make(map[string]*ElasticSource)
var defaultServer string

type ElasticSource struct {
	client *elastic.Client
}

func Init() {
	logs.Debug("|foundation|init|elastic|Init")
	//读取配置
	redissource := ""
	sniff := false
	if redissource != "" {
		InitElasticServer(redissource, &sniff)
	}
}

func InitElasticServer(server string, sniff *bool) bool {
	cfg := &config.Config{URL: server, Sniff: sniff}
	elasticSource, err := NewClient(cfg)
	if err != nil {
		fmt.Sprintf("Elastic Client Init Err:%v", err)
		panic(err)
	}
	elasticSourceMap[server] = elasticSource
	return true
}

func GetElasticDefault() *ElasticSource {
	if len(defaultServer) == 0 {
		for _, s := range elasticSourceMap {
			return s
		}
	}
	return elasticSourceMap[defaultServer]
}

//NewClient -
func NewClient(config *config.Config) (*ElasticSource, error) {
	c, err := elastic.NewClientFromConfig(config)
	if err != nil {
		panic(err)
	}
	return &ElasticSource{c}, err
}

//GetVersion -
func (es *ElasticSource) GetVersion(host string) string {
	esversion, err := es.client.ElasticsearchVersion(host)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	return esversion
}

//CreateTable -Create an index
func (es *ElasticSource) CreateTable(ctx context.Context, tableName string) error {
	exists, err := es.client.IndexExists(tableName).Do(ctx)
	if err != nil {
		return err
	}
	if !exists {
		//创建索引
		_, err := es.client.CreateIndex(tableName).Do(ctx)
		if err != nil {
			return err
		}
	}
	return err
}

//IndexExists -
func (es *ElasticSource) IndexExists(ctx context.Context, tableName string) (bool, error) {
	exists, err := es.client.IndexExists(tableName).Do(ctx)
	return exists, err
}

//Insert -Add a document to the index
func (es *ElasticSource) Insert(ctx context.Context, tableName string, id string, data interface{}) error {
	//Insert
	put1, err := es.client.Index().
		Index(tableName).
		Id(id).
		BodyJson(data).
		Do(ctx)
	fmt.Printf("Indexed elastic %s to index %s, type %s\n", put1.Id, put1.Index, put1.Type)
	return err
}

//UpdateElasticData - update a document to the index
func (es *ElasticSource) Update(ctx context.Context, tableName string, id string, data interface{}) error {
	//Update
	var err error
	_, err = es.client.Update().
		Index(tableName).
		Id(id).
		Doc(data). //map[string]interface{}{"age": 88} or struct
		Do(ctx)
	return err
}

//DeleteElasticData -delete on data
func (es *ElasticSource) Delete(ctx context.Context, tableName string, id string) error {
	var err error
	_, err = es.client.Delete().
		Index(tableName).
		Id(id).
		Do(ctx)
	return err
}

//DropTable - DropTable
func (es *ElasticSource) DropTable(ctx context.Context, tableName string) error {
	var err error
	_, err = es.client.DeleteIndex(tableName).Do(ctx)
	return err
}

//SearchElasticQuery - Search Search with a term query
func (es *ElasticSource) SearchElasticQuery(ctx context.Context, tablename string, key string, value string) (searchResult *elastic.SearchResult, err error) {
	termQuery := elastic.NewTermQuery(key, value)
	searchResult, err = es.client.Search().
		Index(tablename). //search in index "tweets"
		Query(termQuery).
		Sort(key+".keyword", true). //sort by "user" field,ascending
		From(0).Size(10).           //take docuemnt 0-9
		Pretty(true).               //output json
		Do(ctx)
	return
}
