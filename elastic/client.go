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
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	script2 "github.com/lemonyxk/eutils/elastic/script"
	"github.com/lemonyxk/eutils/elastic/types"
	"github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/kitty"
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

type typ int

const (
	aggregation typ = iota + 1
	search
)

func (c *Client) Aggregate(indexes ...string) *Indexes {
	return &Indexes{c, indexes, aggregation}
}

func (c *Client) Search(indexes ...string) *Indexes {
	return &Indexes{c, indexes, search}
}

type Indexes struct {
	client  *Client
	indexes []string
	t       typ
}

func (a *Indexes) Query(query kitty.M) *Response {
	return &Response{a, query}
}

type Response struct {
	*Indexes
	query kitty.M
}

func (a *Response) All(result any) error {

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

	switch a.t {
	case aggregation:
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
	case search:
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
				MaxScore interface{} `json:"max_score"`
				Hits     []struct {
					Index  string          `json:"_index"`
					Id     string          `json:"_id"`
					Score  int             `json:"_score"`
					Source json.RawMessage `json:"_source"`
					Sort   json.RawMessage `json:"sort"`
				} `json:"hits"`
			} `json:"hits"`
		}
		err = json.NewDecoder(res.Body).Decode(&agg)
		if err != nil {
			return err
		}

		var b = bytes.Buffer{}
		b.WriteString(`{"count":`)
		b.WriteString(strconv.Itoa(agg.Hits.Total.Value))
		b.WriteString(`,"list":`)

		if len(agg.Hits.Hits) == 0 {
			b.WriteString("null")
			b.WriteString("}")
			return json.Unmarshal(b.Bytes(), result)
		}

		for i, hit := range agg.Hits.Hits {
			if i == 0 {
				b.WriteByte('[')
			} else {
				b.WriteByte(',')
			}
			b.Write(hit.Source)
			if i == len(agg.Hits.Hits)-1 {
				b.WriteByte(']')
			}
		}

		if len(agg.Hits.Hits) > 0 {
			b.WriteString(`,"sort":`)
			b.Write(agg.Hits.Hits[len(agg.Hits.Hits)-1].Sort)
		}

		b.WriteString("}")

		return json.Unmarshal(b.Bytes(), result)
	default:
		return errors.New("unknown search type")
	}
}
