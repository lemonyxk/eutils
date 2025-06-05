package elastic

import (
	"github.com/lemonyxk/eutils/mapping"
	"github.com/lemonyxk/kitty/kitty"
)

var DynamicTemplates = []kitty.M{
	{
		"string_as_keyword": kitty.M{
			"match_mapping_type": "string",
			"mapping": kitty.M{
				"type":         "keyword",
				"ignore_above": 256,
			},
		},
	},
}

func Ets() *mapping.Mapping {
	var ets = mapping.New()
	ets.DefaultKeyword(false)
	ets.IgnoreNil(false)
	ets.WithTag(false)
	ets.TextAsKeyword(true)
	return ets
}

func MakeMapping[T any]() map[string]any {
	var t T
	return Ets().GenerateMapping(t)
}

func MakeDynamicTemplate[T any](d bool) kitty.M {
	var m = MakeMapping[T]()
	m["dynamic_templates"] = DynamicTemplates
	m["dynamic"] = d
	return m
}

type Model[T Elastic] struct {
	client *Client
	config EsConfig
	t      T
}

type EsConfig struct {
	Format   string
	Prefix   string
	Date     string
	Settings map[string]any
	Mappings map[string]any
}

type Identity interface {
	Timestamp() int64
	String() string
	Empty() bool
}

type ID interface {
	ElasticID() Identity
}

type Config interface {
	Config() EsConfig
}

type Marshaler interface {
	Marshal() ([]byte, error)
}

type Elastic interface {
	ID
	Config
	Marshaler
}
