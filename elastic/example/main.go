/**
* @program: elasticsql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 18:50
**/

package main

import (
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/lemonyxk/eutils/elastic"
	"github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/kitty"
	"time"
)

type ID string

func (a ID) String() string {
	return string(a)
}

func (a ID) Empty() bool {
	return a == ""
}

func (a ID) Timestamp() int64 {
	return 0
}

type Trend struct {
	ID         ID      `json:"id" bson:"_id"`
	PackageID  int     `json:"package_id" bson:"package_id"`
	Counter    int64   `json:"counter" bson:"counter"`
	Value      float64 `json:"value" bson:"value"`
	Remark     string  `json:"remark" bson:"remark"`
	ForTime    int64   `json:"for_time" bson:"for_time"`
	CreateTime int64   `json:"create_time" bson:"create_time" index:"create_time_1"`
}

func (a *Trend) Config() elastic.EsConfig {
	return elastic.EsConfig{
		Format:   "2006",
		Prefix:   "test_1",
		Date:     time.Now().Format("2006"),
		Mappings: elastic.MakeDynamicTemplate[Trend](),
		Settings: kitty.M{
			"index": kitty.M{
				"refresh_interval": "1s",
				"sort": kitty.M{
					"field": []string{"create_time"},
					"order": []string{"desc"},
				},
			},
		},
	}
}

func (a *Trend) Empty() bool {
	return a == nil || a.ID == ""
}

func (a *Trend) ElasticID() elastic.Identity {
	return a.ID
}

func (a *Trend) Identity() string {
	return a.ID.String()
}

func (a *Trend) Marshal() ([]byte, error) {
	type Alias Trend
	var alias = Alias(*a)
	return json.Marshal(alias)
}

func main() {

	//var sql = "select a.b from a where SHOULD(`abc#2` BETWEEN 1 and 2, efg = 2) and `id#10` = 1 and `title#3` is not null and ((name = 'a' or name = 'b') and SHOULD(x = 1 , xx = 2) or c=2 and (age = 1 or age = 2)) and title like `1%1` order by id desc limit 10, 20"

	var t = &Trend{
		ID:         "1",
		PackageID:  1,
		Counter:    1,
		Value:      1,
		Remark:     "1",
		ForTime:    1,
		CreateTime: 1,
	}

	var m = elastic.NewModel[*Trend](&elasticsearch.Client{})
	println(
		elastic.NewBulk(m).Create(t).Upsert(t, elastic.Params{
			Set: kitty.M{"name": "set"},
			Inc: kitty.M{"counter": 1},
		}).String(),
	)

}
