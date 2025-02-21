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
	"time"
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

	//var post = Account{
	//	//TestMap: map[string]interface{}{
	//	//	"test": "test",
	//	//	"test2": 1,
	//	//},
	//	Property: &Property{
	//		Manager: &Manager{
	//			Objects: Objects{
	//				1: &Object{
	//					Type: &Type{
	//						Code: 1,
	//					},
	//				},
	//			},
	//		},
	//		Extends: &Extend{
	//			Objects: Objects{
	//				1: &Object{
	//					Type: &Type{
	//						Code: 1,
	//					},
	//				},
	//			},
	//		},
	//	},
	//	Objects: Objects{
	//		2: &Object{
	//			Type: &Type{
	//				Code: 1,
	//			},
	//		},
	//	},
	//}
	//
	//_ = post
	//
	var ets = mapping.New()
	ets.DefaultKeyword(false)
	ets.IgnoreNil(false)
	ets.WithTag(false)
	ets.TextAsKeyword(true)
	//
	//var mapping = ets.GenerateMapping(post)
	//bts, _ := json.Marshal(mapping)
	//f, err := os.OpenFile(`test.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//f.Write(bts)
	//
	//bts, err = json.Marshal(post)
	//if err != nil {
	//	panic(err)
	//}
	//f1, err := os.OpenFile(`test2.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//defer f1.Close()
	//f1.Write(bts)
	//
	//var a Account
	//err = json.Unmarshal(bts, &a)
	//if err != nil {
	//	panic(err)
	//}
	//
	//log.Printf("%+v", a)

	var runtime = Runtime{}

	mapping := ets.GenerateMapping(runtime)

	bts, _ := json.Marshal(mapping)

	println(string(bts))

	//f, err = os.OpenFile(`test3.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer f.Close()
	//f.Write(bts)
}

type NumericDate struct {
	time.Time
}

type Runtime struct {
	Method   string       `json:"method" bson:"method"`
	Path     string       `json:"path" bson:"path"`
	IP       string       `json:"ip" bson:"ip"`
	Time     *NumericDate `json:"time" bson:"time" es:"type:date"`
	Params   any          `json:"params,omitempty" bson:"params,omitempty" es:"type:flattened"`
	Response any          `json:"response,omitempty" bson:"response,omitempty" es:"type:flattened"`
}
