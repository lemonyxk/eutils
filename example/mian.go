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
	"log"
	"os"
	"time"

	"github.com/lemonyxk/eutils"
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

func main() {

	var post = Post{
		ID:           "",
		Index:        "",
		UserID:       "",
		Type:         0,
		CategoriesID: []string{"a"},
		Tags:         []string{"a"},
		Title:        "",
		Desc:         "",
		Content:      "",
		Covers: []MediaWithSize{
			{
				Media: Media{
					Resource: Resource{
						Path:        "",
						Name:        "",
						StorageType: 0,
					},
					File: File{
						Name: "",
						Size: 0,
						Mime: "",
					},
					Link: Link{
						Domain: "",
						Params: "",
					},
				},
				Size: Size{},
			},
		},
		Videos: []Media{
			{
				Resource: Resource{
					Path:        "",
					Name:        "",
					StorageType: 0,
				},
				File: File{
					Name: "",
					Size: 0,
					Mime: "",
				},
				Link: Link{
					Domain: "",
					Params: "",
				},
			},
		},
		Images: []Media{
			{
				Resource: Resource{
					Path:        "",
					Name:        "",
					StorageType: 0,
				},
				File: File{
					Name: "",
					Size: 0,
					Mime: "",
				},
				Link: Link{
					Domain: "",
					Params: "",
				},
			},
		},
		Audios: []Media{
			{
				Resource: Resource{
					Path:        "",
					Name:        "",
					StorageType: 0,
				},
				File: File{
					Name: "",
					Size: 0,
					Mime: "",
				},
				Link: Link{
					Domain: "",
					Params: "",
				},
			},
		},
		CreateTime: 0,
		UpdateTime: 0,
		Status:     0,
		Record:     PostRecord{},
		User:       nil,
		Map: map[string]interface{}{
			"test":  User{},
			"test1": "hello",
		},
		Slice: []any{1},
	}

	_ = post

	var em = eutils.NewMapping()
	em.DefaultKeyword(true)
	em.WithTag(false)
	em.IgnoreNil(false)

	var start = time.Now()
	var mapping = em.GenerateMapping(post)
	log.Println(time.Since(start))
	bts, _ := json.MarshalIndent(mapping, "", "    ")
	println(string(bts))
	f, err := os.OpenFile(`/Users/lemo/www/test.json`, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.Write(bts)
}
