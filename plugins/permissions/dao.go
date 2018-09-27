package permissions

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/appbaseio-confidential/arc/internal/types/op"
	"github.com/appbaseio-confidential/arc/internal/types/permission"
	"github.com/olivere/elastic"
)

type elasticsearch struct {
	url       string
	indexName string
	typeName  string
	mapping   string
	client    *elastic.Client
}

func NewES(url, indexName, typeName, mapping string) (*elasticsearch, error) {
	opts := []elastic.ClientOptionFunc{
		elastic.SetURL(url),
		elastic.SetSniff(false),
	}
	ctx := context.Background()

	// Initialize the client
	client, err := elastic.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("%s: error while initializing elastic client: %v\n", logTag, err)
	}
	es := &elasticsearch{url, indexName, typeName, mapping, client}

	// Check if the meta index already exists
	exists, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: error while checking if index already exists: %v\n", logTag, err)
	}
	if exists {
		log.Printf("%s index named '%s' already exists, skipping...", logTag, indexName)
		return es, nil
	}

	// Create a new meta index
	_, err = client.CreateIndex(indexName).Body(mapping).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: error while creating index named %s: %v\n", logTag, indexName, err)
	}

	log.Printf("%s successfully created index named '%s'", logTag, indexName)
	return es, nil
}

func (es *elasticsearch) getRawPermissions(username string) ([]byte, error) {
	resp, err := es.client.Get().
		Index(es.indexName).
		Type(es.typeName).
		Id(username).
		FetchSource(true).
		Do(context.Background())
	if err != nil {
		return nil, err
	}

	raw, _ := json.Marshal(resp)
	log.Printf("%s: es_response: %v", logTag, raw)

	src, err := resp.Source.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return src, nil
}

func (es *elasticsearch) putPermission(p permission.Permission) (bool, error) {
	resp, err := es.client.Index().
		Index(es.indexName).
		Type(es.typeName).
		Id(p.UserName).
		BodyJson(p).
		Do(context.Background())
	if err != nil {
		return false, err
	}

	raw, _ := json.Marshal(resp)
	log.Printf("%s: es_response: %s\n", logTag, raw)

	return true, nil
}

func (es *elasticsearch) patchPermission(username string, p permission.Permission) (bool, error) {
	fields := make(map[string]interface{})
	if p.ACL != nil {
		fields["acl"] = p.ACL
	}
	if p.Op != op.Noop {
		fields["op"] = p.Op
	}
	if p.Indices != nil {
		fields["indices"] = p.Indices
	}

	resp, err := es.client.Update().
		Index(es.indexName).
		Type(es.typeName).
		Id(username).
		Doc(fields).
		Do(context.Background())
	if err != nil {
		return false, nil
	}

	raw, _ := json.Marshal(resp)
	log.Printf("%s: es_response: %s\n", logTag, raw)

	return true, nil
}

func (es *elasticsearch) deletePermission(userId string) (bool, error) {
	resp, err := es.client.Delete().
		Index(es.indexName).
		Type(es.typeName).
		Id(userId).
		Do(context.Background())
	if err != nil {
		return false, err
	}

	raw, _ := json.Marshal(resp)
	log.Printf("%s: es_response: %s\n", logTag, raw)

	return true, nil
}