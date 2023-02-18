package backend

import (
	"context"
	"fmt"

	"around/constants"
	"around/util"

	"github.com/olivere/elastic/v7"
)

var (
	ESBackend *ElasticsearchBackend // used for instantiation
)

type ElasticsearchBackend struct { // =DAO in java
	client *elastic.Client // same as sessionFactory in java, used to connect db
}

func InitElasticsearchBackend(config *util.ElasticsearchInfo) {
	// new connection client
	// NewClient(es url, username and password)
	client, err := elastic.NewClient(
        elastic.SetURL(config.Address),
        elastic.SetBasicAuth(config.Username, config.Password))
	// trace exception
	if err != nil {
		panic(err)
	}

	// store post message
	// first check if this index exist or not
	exists, err := client.IndexExists(constants.POST_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	// create mapping, what the table looks like
	if !exists {
		// text means partially equal in search
		// keyword means totally equal in search
		// index: defaultl is true, setup index if frequently use that property to search
		mapping := `{
            "mappings": {
                "properties": {
                    "id":       { "type": "keyword" },
                    "user":     { "type": "keyword" },
                    "message":  { "type": "text" },
                    "url":      { "type": "keyword", "index": false },
                    "type":     { "type": "keyword", "index": false }
                }
            }
        }`
		// create table based on this mapping schema
		_, err := client.CreateIndex(constants.POST_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}

	exists, err = client.IndexExists(constants.USER_INDEX).Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exists {
		mapping := `{
                        "mappings": {
                                "properties": {
                                        "username": {"type": "keyword"},
                                        "password": {"type": "keyword"},
                                        "age":      {"type": "long", "index": false},
                                        "gender":   {"type": "keyword", "index": false}
                                }
                        }
                }`
		_, err = client.CreateIndex(constants.USER_INDEX).Body(mapping).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
	fmt.Println("Indexes are created.")

	// assign the client to package variable
	ESBackend = &ElasticsearchBackend{client: client}
}

func (backend *ElasticsearchBackend) ReadFromES(query elastic.Query, index string) (*elastic.SearchResult, error) {
	searchResult, err := backend.client.Search().
		Index(index). // from in sql
		Query(query). // where in sql
		Pretty(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	return searchResult, nil
}

func (backend *ElasticsearchBackend) SaveToES(i interface{}, index string, id string) error { // interface{} same as java Object, mean take any type as input
	_, err := backend.client.Index().
		Index(index).
		Id(id).
		BodyJson(i).
		Do(context.Background())
	return err
}

func (backend *ElasticsearchBackend) DeleteFromES(query elastic.Query, index string) error {
	_, err := backend.client.DeleteByQuery().
		Index(index).
		Query(query).
		Pretty(true).
		Do(context.Background())

	return err
}
