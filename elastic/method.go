package elastic

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/lemonyxk/eutils/elastic/types"
	"github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/kitty"
	"github.com/lemonyxk/utils/slice"
	"strings"
	"time"
)

type Result[T Elastic] struct {
	Count int   `json:"count"`
	List  []T   `json:"list"`
	Sort  []any `json:"sort,omitempty"`
}

func (r *Result[T]) First() T {
	return slice.Any(r.List).First()
}

func (q *Q) Empty() bool {
	return q.Field == "" || q.Op == "" || kitty.IsNil(q.Value)
}

func NewModel[T Elastic](client *Client) *Model[T] {
	var t T
	var config = (t).Config()

	var model = &Model[T]{
		client: client,
		config: config,
		t:      t,
	}

	return model
}

func (m *Model[T]) IndexName(id Identity) string {
	var date = time.Unix(id.Timestamp(), 0).Format(m.config.Format)
	return m.config.Prefix + "-" + date
}

func (m *Model[T]) UpdateIndexSettings(settings kitty.M) *Model[T] {
	var body, err = json.Marshal(settings)
	if err != nil {
		panic(err)
	}

	var allowNoIndices = true
	var req = esapi.IndicesPutSettingsRequest{
		Index:          []string{m.config.Prefix + "*"},
		Body:           bytes.NewReader(body),
		AllowNoIndices: &allowNoIndices,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	return m
}

func (m *Model[T]) UpdateIndexMapping(mappings kitty.M) *Model[T] {
	var body, err = json.Marshal(mappings)
	if err != nil {
		panic(err)
	}

	var allowNoIndices = true
	var req = esapi.IndicesPutMappingRequest{
		Index:          []string{m.config.Prefix + "*"},
		Body:           bytes.NewReader(body),
		AllowNoIndices: &allowNoIndices,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	return m
}

// MakeIndexTemplate
// https://www.elastic.co/guide/en/elasticsearch/reference/current/simulate-multi-component-templates.html
func (m *Model[T]) MakeIndexTemplate() *Model[T] {

	var config = kitty.M{
		"index_patterns": []string{m.config.Prefix + "*"},
		//"translog": kitty.M{
		//	"durability": "request", // default is request, can be changed to async
		//},
		"template": kitty.M{
			"settings": m.config.Settings,
			"mappings": m.config.Mappings,
		},
	}

	var body, err = json.Marshal(config)
	if err != nil {
		panic(err)
	}

	var req = esapi.IndicesPutIndexTemplateRequest{
		Name: m.config.Prefix,
		Body: bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	println(res.String())

	return m
}

//

// MakeMapping
// https://www.elastic.co/guide/en/elasticsearch/reference/current/index-modules-translog.html
// how to make es do not lose data
func (m *Model[T]) MakeMapping(date string) *Model[T] {

	var index = m.config.Prefix + "-" + date

	var config = kitty.M{
		"settings": m.config.Settings,
		"mappings": m.config.Mappings,
	}

	var body, err = json.Marshal(config)
	if err != nil {
		panic(err)
	}

	var req = esapi.IndicesCreateRequest{
		Index: index,
		Body:  bytes.NewReader(body),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() && !strings.Contains(res.String(), "already exists") {
		panic(res.String())
	}

	mts, err := json.Marshal(m.config.Mappings)
	if err != nil {
		panic(err)
	}

	var sr = esapi.IndicesPutMappingRequest{
		Index: []string{index},
		Body:  bytes.NewReader(mts),
	}

	res, err = sr.Do(context.Background(), m.client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	return m
}

// Get returns a document from the index.
func (m *Model[T]) Get(id Identity) (T, error) {
	var t T

	var req = esapi.GetRequest{
		Index:      m.IndexName(id),
		DocumentID: id.String(),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return t, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		if res.StatusCode == 404 {
			return t, nil
		}
		return t, errors.New(res.String())
	}

	var data struct {
		Source T `json:"_source"`
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return t, err
	}

	return data.Source, nil
}

func (m *Model[T]) Gets(ids ...Identity) ([]T, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	var list []kitty.M

	for i := 0; i < len(ids); i++ {
		list = append(list, kitty.M{"_index": m.IndexName(ids[i]), "_id": ids[i]})
	}

	var body = kitty.M{"docs": list}

	var bts, err = json.Marshal(body)
	if err != nil {
		return nil, err
	}

	var req = esapi.MgetRequest{
		Body: bytes.NewReader(bts),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var data struct {
		Docs []struct {
			Found  bool `json:"found"`
			Source T    `json:"_source"`
		}
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var result []T
	for i := 0; i < len(data.Docs); i++ {
		if data.Docs[i].Found {
			result = append(result, data.Docs[i].Source)
		}
	}

	return result, nil
}

func (m *Model[T]) Query(query Query) (*Result[T], error) {

	var result = &Result[T]{}

	if query.Limit == 0 {
		query.Limit = 1000
	}

	var sql SQL
	var err error
	if query.SQL != "" {
		sql = query.SQL
	} else {
		sql, err = MakeSQL(QueryTable{Name: m.config.Prefix, Query: query})
		if err != nil {
			return result, err
		}
	}

	var dsl string
	var table string
	if len(query.DSL) > 0 {
		str, err := json.Marshal(query.DSL)
		if err != nil {
			return result, err
		}
		dsl = string(str)
		table = m.config.Prefix
	} else {
		dsl, table, err = ConvertSQLToDSL(sql)
		if err != nil {
			return result, err
		}
	}

	if table != m.config.Prefix {
		return result, errors.New("table name error")
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var indexes []string
	for i := 0; i < len(query.Indexes); i++ {
		indexes = append(indexes, table+"-"+query.Indexes[i])
	}
	if len(indexes) == 0 {
		indexes = []string{table + "*"}
	}

	var req = esapi.SearchRequest{
		Index: indexes,
		Body:  strings.NewReader(dsl),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return result, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return result, errors.New(res.String())
	}

	var source struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`

			Hits []struct {
				Source T     `json:"_source"`
				Sort   []any `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&source)
	if err != nil {
		return result, err
	}

	var list []T

	for i := 0; i < len(source.Hits.Hits); i++ {
		list = append(list, source.Hits.Hits[i].Source)
	}

	result.Count = source.Hits.Total.Value
	result.List = list

	if len(source.Hits.Hits) > 0 {
		result.Sort = source.Hits.Hits[len(source.Hits.Hits)-1].Sort
	}

	return result, nil
}

func (m *Model[T]) Queries(queries ...Query) ([]*Result[T], error) {

	if len(queries) == 0 {
		return nil, errors.New("queries is empty")
	}

	var results []*Result[T]

	var indexes []string

	var queryStr string
	for i := 0; i < len(queries); i++ {
		var sql SQL
		var err error
		if queries[i].SQL != "" {
			sql = queries[i].SQL
		} else {
			sql, err = MakeSQL(QueryTable{Name: m.config.Prefix, Query: queries[i]})
			if err != nil {
				return results, err
			}
		}

		var dsl string
		var table string
		if len(queries[i].DSL) > 0 {
			str, err := json.Marshal(queries[i].DSL)
			if err != nil {
				return results, err
			}
			dsl = string(str)
			table = m.config.Prefix
		} else {
			dsl, table, err = ConvertSQLToDSL(sql)
			if err != nil {
				return results, err
			}
		}

		if table != m.config.Prefix {
			return results, errors.New("table name error")
		}

		queryStr += `{}` + "\n" + dsl + "\n"

		for j := 0; j < len(queries[i].Indexes); j++ {
			indexes = append(indexes, table+"-"+queries[i].Indexes[j])
		}
		if len(indexes) == 0 {
			indexes = []string{table + "*"}
		}
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", queryStr, time.Since(now))
	}()

	indexes = slice.Compare(indexes).Unique()

	var req = esapi.MsearchRequest{
		Index: indexes,
		Body:  strings.NewReader(queryStr),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return results, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return results, errors.New(res.String())
	}

	var source struct {
		Responses []struct {
			Hits struct {
				Total struct {
					Value int `json:"value"`
				} `json:"total"`

				Hits []struct {
					Source T     `json:"_source"`
					Sort   []any `json:"sort"`
				} `json:"hits"`
			} `json:"hits"`
		} `json:"responses"`
	}

	err = json.NewDecoder(res.Body).Decode(&source)
	if err != nil {
		return results, err
	}

	for i := 0; i < len(source.Responses); i++ {
		var response = source.Responses[i]

		var result = &Result[T]{}

		var list []T

		for j := 0; j < len(response.Hits.Hits); j++ {
			list = append(list, response.Hits.Hits[j].Source)
		}

		result.Count = response.Hits.Total.Value
		result.List = list

		if len(response.Hits.Hits) > 0 {
			result.Sort = response.Hits.Hits[len(response.Hits.Hits)-1].Sort
		}

		results = append(results, result)
	}

	return results, nil
}

func (m *Model[T]) Count(query Query) (int, error) {

	var result = 0

	query.Limit = 0
	query.Page = 1

	var sql SQL
	var err error
	if query.SQL != "" {
		sql = query.SQL
	} else {
		sql, err = MakeSQL(QueryTable{Name: m.config.Prefix, Query: query})
		if err != nil {
			return result, err
		}
	}

	var dsl string
	var table string
	if len(query.DSL) > 0 {
		str, err := json.Marshal(query.DSL)
		if err != nil {
			return result, err
		}
		dsl = string(str)
		table = m.config.Prefix
	} else {
		dsl, table, err = ConvertSQLToDSL(sql)
		if err != nil {
			return result, err
		}
	}

	if table != m.config.Prefix {
		return result, errors.New("table name error")
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var indexes []string
	for i := 0; i < len(query.Indexes); i++ {
		indexes = append(indexes, table+"-"+query.Indexes[i])
	}
	if len(indexes) == 0 {
		indexes = []string{table + "*"}
	}

	var req = esapi.SearchRequest{
		Index:          indexes,
		Body:           strings.NewReader(dsl),
		TrackTotalHits: true,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return result, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return result, errors.New(res.String())
	}

	var source struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`

			Hits []struct {
				Source T `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&source)
	if err != nil {
		return result, err
	}

	result = source.Hits.Total.Value

	return result, nil
}

func (m *Model[T]) Find(query kitty.M) (*Result[T], error) {
	var dsl, err = json.Marshal(query)
	if err != nil {
		return nil, err
	}

	return m.DSL(string(dsl))
}

func (m *Model[T]) SQL(sql SQL) (*Result[T], error) {
	dsl, _, err := ConvertSQLToDSL(sql)
	if err != nil {
		return nil, err
	}

	return m.DSL(dsl)
}

func (m *Model[T]) DSL(dsl string) (*Result[T], error) {
	var result = &Result[T]{}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var req = esapi.SearchRequest{
		Index: []string{m.config.Prefix + "*"},
		Body:  strings.NewReader(dsl),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return result, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return result, errors.New(res.String())
	}

	var source struct {
		Hits struct {
			Total struct {
				Value int `json:"value"`
			} `json:"total"`

			Hits []struct {
				Source T     `json:"_source"`
				Sort   []any `json:"sort"`
			} `json:"hits"`
		} `json:"hits"`
	}

	err = json.NewDecoder(res.Body).Decode(&source)
	if err != nil {
		return result, err
	}

	var list []T

	for i := 0; i < len(source.Hits.Hits); i++ {
		list = append(list, source.Hits.Hits[i].Source)
	}

	result.Count = source.Hits.Total.Value
	result.List = list

	if len(source.Hits.Hits) > 0 {
		result.Sort = source.Hits.Hits[len(source.Hits.Hits)-1].Sort
	}

	return result, nil
}

func (m *Model[T]) Search(search kitty.M) ([]byte, error) {

	var bts, err = json.Marshal(search)
	if err != nil {
		return nil, err
	}

	var dsl = string(bts)

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var req = esapi.SearchRequest{
		Index: []string{m.config.Prefix + "*"},
		Body:  strings.NewReader(dsl),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (m *Model[T]) Searches(searches []kitty.M) ([]byte, error) {

	var dsl string
	for i := 0; i < len(searches); i++ {
		var bts, err = json.Marshal(searches[i])
		if err != nil {
			return nil, err
		}
		dsl += `{}` + "\n" + string(bts) + "\n"
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", dsl, time.Since(now))
	}()

	var req = esapi.MsearchRequest{
		Index: []string{m.config.Prefix + "*"},
		Body:  strings.NewReader(dsl),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var buffer bytes.Buffer
	_, err = buffer.ReadFrom(res.Body)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (m *Model[T]) Aggregate(query kitty.M) *AggregationResponse {
	return m.client.Aggregate(m.config.Prefix + "*").Query(query)
}

// Indexes insert a document into the index.
// or update a document if it already exists.
func (m *Model[T]) Indexes(insert ...T) (*types.MultiIndexResponse, error) {

	var resBuf = bytes.NewBuffer(nil)
	for i := 0; i < len(insert); i++ {
		var t = insert[i]
		resBuf.WriteString(fmt.Sprintf(`{"index":{"_id":"%s","_index":"%s"}}`, t.ElasticID(), m.IndexName(t.ElasticID())))
		resBuf.WriteByte('\n')
		var b, err = t.Marshal()
		if err != nil {
			return nil, err
		}
		resBuf.Write(b)
		resBuf.WriteByte('\n')
	}

	var req = esapi.BulkRequest{
		Body:    resBuf,
		Timeout: time.Second * 30,
	}

	res, err := req.Do(context.Background(), m.client)
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

func (m *Model[T]) Create(insert T) (*types.UpdateResponse, error) {

	var id = insert.ElasticID()

	if id == nil {
		return nil, errors.New("id is empty")
	}

	var date = time.Unix(insert.ElasticID().Timestamp(), 0).Format(m.config.Format)

	bts, err := insert.Marshal()
	if err != nil {
		return nil, err
	}

	var req = esapi.IndexRequest{
		Index:      m.config.Prefix + "-" + date,
		DocumentID: id.String(),
		Body:       bytes.NewReader(bts),
		OpType:     "create",
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Index(insert T) (*types.UpdateResponse, error) {

	var id = insert.ElasticID()

	if id == nil {
		return nil, errors.New("id is empty")
	}

	var date = time.Unix(insert.ElasticID().Timestamp(), 0).Format(m.config.Format)

	bts, err := insert.Marshal()
	if err != nil {
		return nil, err
	}

	var req = esapi.IndexRequest{
		Index:      m.config.Prefix + "-" + date,
		DocumentID: id.String(),
		Body:       bytes.NewReader(bts),
		OpType:     "index",
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Patch(id Identity, t any) (*types.UpdateResponse, error) {

	var date = time.Unix(id.Timestamp(), 0).Format(m.config.Format)

	var body string

	switch t.(type) {
	case T:
		var bts, err = t.(T).Marshal()
		if err != nil {
			return nil, err
		}
		body = fmt.Sprintf(`{"doc": %s}`, string(bts))
	default:
		var bts, err = json.Marshal(t)
		if err != nil {
			return nil, err
		}
		body = fmt.Sprintf(`{"doc": %s}`, string(bts))
	}

	var now = time.Now()
	defer func() {
		fmt.Println("patch:", string(body), time.Since(now))
	}()

	var retryOnConflict = 99999

	var req = esapi.UpdateRequest{
		Index:           m.config.Prefix + "-" + date,
		DocumentID:      id.String(),
		Body:            strings.NewReader(body),
		RetryOnConflict: &retryOnConflict,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		if res.StatusCode == 404 {
			return &types.UpdateResponse{}, nil
		}
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Update(id Identity, params Params) (*types.UpdateResponse, error) {
	var date = time.Unix(id.Timestamp(), 0).Format(m.config.Format)

	var script = Script{
		ID:     "update",
		Params: params,
	}

	scriptBts, err := json.Marshal(script)
	if err != nil {
		return nil, err
	}

	var body = fmt.Sprintf(`{"script": %s}`, string(scriptBts))

	var now = time.Now()
	defer func() {
		fmt.Println("search:", string(body), time.Since(now))
	}()

	var retryOnConflict = 99999

	var req = esapi.UpdateRequest{
		Index:           m.config.Prefix + "-" + date,
		DocumentID:      id.String(),
		Body:            strings.NewReader(body),
		RetryOnConflict: &retryOnConflict,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Upsert(id Identity, t T, params Params) (*types.UpdateResponse, error) {

	var date = time.Unix(id.Timestamp(), 0).Format(m.config.Format)

	var upsertBts, err = t.Marshal()
	if err != nil {
		return nil, err
	}

	var script = Script{
		ID:     "update",
		Params: params,
	}

	scriptBts, err := json.Marshal(script)
	if err != nil {
		return nil, err
	}

	var body = fmt.Sprintf(`{"script": %s, "upsert": %s}`, string(scriptBts), string(upsertBts))

	var now = time.Now()
	defer func() {
		fmt.Println("search:", string(body), time.Since(now))
	}()

	var retryOnConflict = 99999

	var req = esapi.UpdateRequest{
		Index:           m.config.Prefix + "-" + date,
		DocumentID:      id.String(),
		Body:            strings.NewReader(body),
		RetryOnConflict: &retryOnConflict,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Modify(query Query, params Params) (*types.UpdateByQueryResponse, error) {

	var sql SQL
	var err error
	if query.SQL != "" {
		sql = query.SQL
	} else {
		sql, err = MakeSQL(QueryTable{Name: m.config.Prefix, Query: query})
		if err != nil {
			return nil, err
		}
	}

	var dsl string
	var table string
	if len(query.DSL) > 0 {
		str, err := json.Marshal(query.DSL)
		if err != nil {
			return nil, err
		}
		dsl = string(str)
		table = m.config.Prefix
	} else {
		dsl, table, err = ConvertSQLToDSL(sql)
		if err != nil {
			return nil, err
		}
	}

	if table != m.config.Prefix {
		return nil, errors.New("table name error")
	}

	var script = Script{
		ID:     "update",
		Params: params,
	}

	scriptBts, err := json.Marshal(script)
	if err != nil {
		return nil, err
	}

	var qs struct {
		Query json.RawMessage `json:"query"`
	}

	err = json.Unmarshal([]byte(dsl), &qs)
	if err != nil {
		return nil, err
	}

	var uq = `{` + fmt.Sprintf(`"script": %s, "query": %s`, string(scriptBts), string(qs.Query)) + `}`
	if string(qs.Query) == "" {
		uq = `{` + fmt.Sprintf(`"script": %s`, string(scriptBts)) + `}`
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", uq, time.Since(now))
	}()

	var indexes []string
	for i := 0; i < len(query.Indexes); i++ {
		indexes = append(indexes, table+"-"+query.Indexes[i])
	}
	if len(indexes) == 0 {
		indexes = []string{table + "*"}
	}

	var waitForCompletion = true
	var req = esapi.UpdateByQueryRequest{
		Index:             indexes,
		Body:              strings.NewReader(uq),
		Conflicts:         "proceed",
		WaitForCompletion: &waitForCompletion,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.UpdateByQueryResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

// Delete deletes a document from the index.
func (m *Model[T]) Delete(id Identity) (*types.UpdateResponse, error) {

	var date = time.Unix(id.Timestamp(), 0).Format(m.config.Format)

	var req = esapi.DeleteRequest{
		Index:      m.config.Prefix + "-" + date,
		DocumentID: id.String(),
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	var esResponse types.UpdateResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	if res.IsError() {
		if res.StatusCode == 404 {
			return &esResponse, nil
		}
		return nil, errors.New(res.String())
	}

	return &esResponse, nil
}

func (m *Model[T]) Remove(query Query) (*types.DeleteByQueryResponse, error) {

	var sql SQL
	var err error
	if query.SQL != "" {
		sql = query.SQL
	} else {
		sql, err = MakeSQL(QueryTable{Name: m.config.Prefix, Query: query})
		if err != nil {
			return nil, err
		}
	}

	var dsl string
	var table string
	if len(query.DSL) > 0 {
		str, err := json.Marshal(query.DSL)
		if err != nil {
			return nil, err
		}
		dsl = string(str)
		table = m.config.Prefix
	} else {
		dsl, table, err = ConvertSQLToDSL(sql)
		if err != nil {
			return nil, err
		}
	}

	if table != m.config.Prefix {
		return nil, errors.New("table name error")
	}

	var qs struct {
		Query json.RawMessage `json:"query"`
	}

	err = json.Unmarshal([]byte(dsl), &qs)
	if err != nil {
		return nil, err
	}

	var uq = `{` + fmt.Sprintf(`"query": %s`, string(qs.Query)) + `}`

	var now = time.Now()
	defer func() {
		fmt.Println("search:", uq, time.Since(now))
	}()

	var indexes []string
	for i := 0; i < len(query.Indexes); i++ {
		indexes = append(indexes, table+"-"+query.Indexes[i])
	}
	if len(indexes) == 0 {
		indexes = []string{table + "*"}
	}

	var waitForCompletion = true
	var req = esapi.DeleteByQueryRequest{
		Index:             indexes,
		Body:              strings.NewReader(uq),
		Conflicts:         "proceed",
		WaitForCompletion: &waitForCompletion,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return nil, err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var esResponse types.DeleteByQueryResponse
	err = json.NewDecoder(res.Body).Decode(&esResponse)
	if err != nil {
		return nil, err
	}

	return &esResponse, nil
}

func (m *Model[T]) Reindex(from, to string, script ...string) (string, error) {

	var scriptStr string
	for i := 0; i < len(script); i++ {
		scriptStr += script[i] + ";"
	}

	var body string

	if scriptStr == "" {
		body = fmt.Sprintf(`{"source":{"index":"%s"}, "dest":{"index":"%s"}}`,
			m.config.Prefix+"-"+from, m.config.Prefix+"-"+to)
	} else {
		body = fmt.Sprintf(`{"source":{"index":"%s"}, "dest":{"index":"%s"}, "script":{"source":"%s"}}`,
			m.config.Prefix+"-"+from, m.config.Prefix+"-"+to, scriptStr)
	}

	var now = time.Now()
	defer func() {
		fmt.Println("search:", body, time.Since(now))
	}()

	var waitForCompletion = true
	var req = esapi.ReindexRequest{
		Body:              strings.NewReader(body),
		WaitForCompletion: &waitForCompletion,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return "", err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		return "", errors.New(res.String())
	}

	return res.String(), nil
}

func (m *Model[T]) Drop(date ...string) (string, error) {

	var indexes []string

	if len(date) == 0 {
		return "", errors.New("date is empty")
	}

	for i := 0; i < len(date); i++ {
		indexes = append(indexes, m.config.Prefix+"-"+date[i])
	}

	var now = time.Now()
	defer func() {
		fmt.Println("drop:", indexes, time.Since(now))
	}()

	var req = esapi.IndicesDeleteRequest{
		Index: indexes,
	}

	res, err := req.Do(context.Background(), m.client)
	if err != nil {
		return "", err
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		if res.StatusCode == 404 {
			return "", nil
		}
		return "", errors.New(res.String())
	}

	return res.String(), nil
}
