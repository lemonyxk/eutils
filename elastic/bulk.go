/**
* @program: eutils
*
* @create: 2025-04-21 17:28
**/

/*

POST _bulk
{ "index" : { "_index" : "test", "_id" : "1" } }
{ "field1" : "value1" }
{ "delete" : { "_index" : "test", "_id" : "2" } }
{ "create" : { "_index" : "test", "_id" : "3" } }
{ "field1" : "value3" }
{ "update" : {"_id" : "1", "_index" : "test"} }
{ "doc" : {"field2" : "value2"} }


POST /_bulk
{ "update" : {"_id" : "1", "_index" : "test_1"} }
{ "doc" : {"field2" : "value21"} ,"upsert":{"aaa":2}}
{ "update" : {"_id" : "2", "_index" : "test_1"} }
{"script": {"id":"update","params": {"$inc":{"trend.download.a123":1,"trend.agree":1,"trend.score.times":1,"trend.score.count":2,"age":1}}}, "upsert": {"name":2,"age":1}}
*/

package elastic

import (
	"fmt"
	"github.com/lemonyxk/kitty/json"
	"strings"
)

type Operation string

const (
	Create Operation = "create"
	Update Operation = "update"
	Index  Operation = "index"
	Delete Operation = "delete"
)

func NewBulk[T Elastic](model *Model[T]) *Bulk[T] {
	return &Bulk[T]{model: model}
}

type Meta struct {
	Index string `json:"_index,omitempty"`
	ID    string `json:"_id,omitempty"`
}

type BulkModel struct {
	Meta     map[Operation]Meta `json:"meta,omitempty"`
	Document string             `json:"document,omitempty"`
}

type BulkModels []*BulkModel

func (d BulkModels) String() string {
	var builder = strings.Builder{}
	for index, data := range d {
		for op, meta := range data.Meta {
			builder.WriteString(fmt.Sprintf(`{"%s":{"_index":"%s","_id":"%s"}}`, op, meta.Index, meta.ID))
		}
		if len(data.Document) > 0 {
			builder.WriteString("\n")
			builder.WriteString(data.Document)
		}
		if index != len(d)-1 {
			builder.WriteString("\n")
		}
	}
	return builder.String()
}

type Bulk[T Elastic] struct {
	model *Model[T]
	list  BulkModels
}

func (d *Bulk[T]) String() string {
	return d.list.String()
}

func (d *Bulk[T]) List() BulkModels {
	return d.list
}

func (d *Bulk[T]) Create(doc Elastic) *Bulk[T] {
	bts, err := doc.Marshal()
	if err != nil {
		panic(err)
	}
	d.list = append(d.list, &BulkModel{
		Document: string(bts),
		Meta: map[Operation]Meta{
			Create: {
				Index: d.model.IndexName(doc.ElasticID()),
				ID:    doc.ElasticID().String(),
			},
		},
	})
	return d
}

func (d *Bulk[T]) Update(doc Elastic) *Bulk[T] {
	bts, err := doc.Marshal()
	if err != nil {
		panic(err)
	}
	d.list = append(d.list, &BulkModel{
		Document: `{"doc":` + string(bts) + `}`,
		Meta: map[Operation]Meta{
			Update: {
				Index: d.model.IndexName(doc.ElasticID()),
				ID:    doc.ElasticID().String(),
			},
		},
	})
	return d
}

func (d *Bulk[T]) Index(doc Elastic) *Bulk[T] {
	bts, err := doc.Marshal()
	if err != nil {
		panic(err)
	}
	d.list = append(d.list, &BulkModel{
		Document: string(bts),
		Meta: map[Operation]Meta{
			Index: {
				Index: d.model.IndexName(doc.ElasticID()),
				ID:    doc.ElasticID().String(),
			},
		},
	})
	return d
}

func (d *Bulk[T]) Delete(id Identity) *Bulk[T] {
	d.list = append(d.list, &BulkModel{
		Meta: map[Operation]Meta{
			Delete: {
				Index: d.model.IndexName(id),
				ID:    id.String(),
			},
		},
	})
	return d
}

func (d *Bulk[T]) Upsert(doc Elastic, params Params) *Bulk[T] {
	bts, err := doc.Marshal()
	if err != nil {
		panic(err)
	}

	var script = Script{
		ID:     "update",
		Params: params,
	}

	scriptBts, err := json.Marshal(script)
	if err != nil {
		panic(err)
	}

	d.list = append(d.list, &BulkModel{
		Document: fmt.Sprintf(`{"script":%s,"upsert":%s}`, string(scriptBts), string(bts)),
		Meta: map[Operation]Meta{
			Update: {
				Index: d.model.IndexName(doc.ElasticID()),
				ID:    doc.ElasticID().String(),
			},
		},
	})
	return d
}

func (d *Bulk[T]) Modify(id Identity, params Params) *Bulk[T] {
	var script = Script{
		ID:     "update",
		Params: params,
	}

	scriptBts, err := json.Marshal(script)
	if err != nil {
		panic(err)
	}

	d.list = append(d.list, &BulkModel{
		Document: fmt.Sprintf(`{"script":%s}`, string(scriptBts)),
		Meta: map[Operation]Meta{
			Update: {
				Index: d.model.IndexName(id),
				ID:    id.String(),
			},
		},
	})
	return d
}
