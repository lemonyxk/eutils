/**
* @program: es
*
* @description:
*
* @author: lemo
*
* @create: 2023-05-21 18:16
**/

package main

import (
	"encoding/json"
	"github.com/lemonyxk/eutils/mapping"
	"os"
)

type Empty interface {
	Empty() bool
}

type Company struct {
	ID            int    `json:"id" es:"index:true"`
	Alias         string `json:"alias" es:"type:keyword,index:true"`
	Name          string `json:"name" es:"type:text"`
	Description   string `json:"description" es:"type:text"`
	EmployeeCount int    `json:"employee_count" es:"index:false"`
	URLs          Empty  `json:"urls"`
}

func (c *Company) Empty() bool {
	return c.URLs == nil
}

func main() {

	var post = Trend{
		//TestMap: map[string]interface{}{
		//	"test": "test",
		//	"test2": 1,
		//},
		Item: Item{
			ID:         "",
			PackageID:  0,
			Name:       TestName{},
			Type:       0,
			Tags:       nil,
			Expire:     Expire{},
			Price:      Price{},
			Images:     nil,
			Context:    Context{},
			CreateTime: 0,
			UpdateTime: 0,
			StartTime:  0,
			EndTime:    0,
			Sort:       0,
			Status:     0,
			Post1: Post1{
				//Name:  "1",
				Link: "1",
				Post2: Post2{
					Name1: "2",
				},
			},
		},
	}

	_ = post

	var ets = mapping.New()
	ets.DefaultKeyword(false)
	ets.IgnoreNil(false)
	ets.WithTag(false)
	ets.TextAsKeyword(true)

	var mapping = ets.GenerateMapping(post)
	bts, _ := json.Marshal(mapping)
	f, err := os.OpenFile(`test.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(bts)

	bts, err = json.Marshal(post)
	if err != nil {
		panic(err)
	}
	f1, err := os.OpenFile(`test2.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	f1.Write(bts)
}
