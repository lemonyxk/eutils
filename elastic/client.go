/**
* @program: eutils
*
* @create: 2025-04-25 23:09
**/

package elastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	script2 "github.com/lemonyxk/eutils/elastic/script"
	"github.com/lemonyxk/eutils/elastic/types"
	"github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/kitty"
	"strings"
	"time"
)

func NewClient(cfg elasticsearch.Config) (*Client, error) {
	var es, err = elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	info, err := es.Info()
	if err != nil {
		return nil, err
	}

	defer func() { _ = info.Body.Close() }()

	println(info.String())

	return &Client{es}, nil
}

type Client struct {
	*elasticsearch.Client
}

func (c *Client) Bulk(models ...*BulkModel) (*types.MultiIndexResponse, error) {

	var req = esapi.BulkRequest{
		Body:    BulkModels(models).Buffer(),
		Timeout: time.Second * 30,
	}

	res, err := req.Do(context.Background(), c)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var multiIndexResponse types.MultiIndexResponse
	err = json.NewDecoder(res.Body).Decode(&multiIndexResponse)
	if err != nil {
		return nil, err
	}

	return &multiIndexResponse, nil
}

func (c *Client) CreateUpdateScript() {

	var script = kitty.M{
		"script": kitty.M{
			"source": script2.UpdateScript,
			"lang":   "painless",
		},
	}

	var bts, err = json.Marshal(script)
	if err != nil {
		panic(err)
	}

	var req = esapi.PutScriptRequest{
		ScriptID: "update",
		Body:     bytes.NewReader(bts),
	}

	res, err := req.Do(context.Background(), c)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	defer func() { _ = res.Body.Close() }()

	println(res.String())
}

func (c *Client) Aggregate(indexes ...string) *Aggregation {
	return &Aggregation{c, indexes}
}

type Aggregation struct {
	client  *Client
	indexes []string
}

func (a *Aggregation) Query(query kitty.M) *AggregationResponse {
	return &AggregationResponse{a.client, a.indexes, query}
}

type AggregationResponse struct {
	client  *Client
	indexes []string
	query   kitty.M
}

func (a *AggregationResponse) All(result any) error {

	var dslBts, err = json.Marshal(a.query)
	if err != nil {
		return err
	}

	var dsl = string(dslBts)

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var req = esapi.SearchRequest{
		Index: a.indexes,
		Body:  strings.NewReader(dsl),
	}

	res, err := req.Do(context.Background(), a.client)
	if err != nil {
		return err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return errors.New(res.String())
	}

	var agg struct {
		Took     int  `json:"took"`
		TimedOut bool `json:"timed_out"`
		Shards   struct {
			Total      int `json:"total"`
			Successful int `json:"successful"`
			Skipped    int `json:"skipped"`
			Failed     int `json:"failed"`
		} `json:"_shards"`
		Hits struct {
			Total struct {
				Value    int    `json:"value"`
				Relation string `json:"relation"`
			} `json:"total"`
			MaxScore interface{}   `json:"max_score"`
			Hits     []interface{} `json:"hits"`
		} `json:"hits"`
		Aggregations json.RawMessage `json:"aggregations"`
	}
	err = json.NewDecoder(res.Body).Decode(&agg)
	if err != nil {
		return err
	}

	return json.Unmarshal(agg.Aggregations, result)
}
