/**
* @program: auth-server
*
* @description:
*
* @author: lemo
*
* @create: 2023-03-02 14:50
**/

package main

import (
	"encoding/json"
)

type Media struct {
	Resource Resource `json:"resource" bson:"resource"`
	File     File     `json:"file" bson:"file"`
	Link     Link     `json:"link" bson:"link"`
}

type Info struct {
	Duration float64 `json:"duration" bson:"duration"`
}

type testAnonymous struct {
}

type AudioWithInfo struct {
	*Media `json:",omitempty" bson:"inline,omitempty"`
	Info   Info `json:"info" bson:"info"`
}

type VideoWithInfo struct {
	*Media   `json:",omitempty" bson:"inline,omitempty"`
	Info     Info           `json:"info" bson:"info"`
	Property *VideoProperty `json:"property,omitempty" bson:"property,omitempty"`
}

type ImageWithSize struct {
	*Media `json:",omitempty" bson:"inline,omitempty"`
	Size   Size `json:"size" bson:"size"`
}

type Size struct {
	Width  int `json:"width" bson:"width"`
	Height int `json:"height" bson:"height"`
}

type Link struct {
	Domain string `json:"domain" bson:"domain"`
	Params string `json:"params" bson:"params"`
}

type Resource struct {
	Path        string `json:"path" bson:"path"`
	Name        string `json:"name" bson:"name"`
	StorageType int    `json:"storage_type" bson:"storage_type"`
}

type File struct {
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	Mime string `json:"mime" bson:"mime"`
}

type PostTrend struct {
	Agree    int `json:"agree" bson:"agree"`
	DisAgree int `json:"disagree" bson:"disagree"`
	Like     int `json:"like" bson:"like"`
	Share    int `json:"share" bson:"share"`
	Collect  int `json:"collect" bson:"collect"`
	Comment  int `json:"comment" bson:"comment"`
	Read     int `json:"read" bson:"read"`
	Play     int `json:"play" bson:"play"`
}

type PostList []*Post

type Magnet struct {
	Name string `json:"name" bson:"name"`
	Link string `json:"link" bson:"link"`
}

type VideoProperty struct {
	Name         string   `json:"name" bson:"name"`
	ActorList    []string `json:"actor_list" bson:"actor_list"`
	DirectorList []string `json:"director_list" bson:"director_list"`
	TypeList     []string `json:"type_list" bson:"type_list"`
	LanguageList []string `json:"language_list" bson:"language_list"`
	Plot         string   `json:"plot" bson:"plot"`
	Score        float64  `json:"score" bson:"score"`
	ScoreCount   int      `json:"score_count" bson:"score_count"`
	ReleaseTime  int64    `json:"release_time" bson:"release_time"`
	ReleaseArea  string   `json:"release_area" bson:"release_area"`
	MagnetList   []Magnet `json:"magnet_list" bson:"magnet_list"`
}

type Post struct {
	ID            string          `json:"id" bson:"_id"`
	EID           int64           `json:"eid" bson:"eid" index:"eid_1,unique"`
	UserID        string          `json:"user_id" bson:"user_id" index:"user_id_1"`
	PackageID     int             `json:"package_id" bson:"package_id"`
	Type          int             `json:"type" bson:"type"`
	CategoriesID  []string        `json:"categories_id" bson:"categories_id" index:"categories_id_1"`
	Tags          []string        `json:"tags" bson:"tags" index:"tags_1" es:"analyzer:ik_max_word"`
	Title         string          `json:"title" bson:"title" index:"title_1" es:"analyzer:ik_max_word"`
	Desc          string          `json:"desc" bson:"desc" es:"analyzer:ik_smart"`
	Content       string          `json:"content" bson:"content" es:"analyzer:ik_smart"`
	Covers        []ImageWithSize `json:"covers" bson:"covers"`
	Images        []ImageWithSize `json:"images" bson:"images"`
	Videos        []VideoWithInfo `json:"videos" bson:"videos"`
	Audios        []AudioWithInfo `json:"audios" bson:"audios"`
	CreateTime    int64           `json:"create_time" bson:"create_time" index:"create_time_1"`
	UpdateTime    int64           `json:"update_time" bson:"update_time" index:"update_time_1"`
	Status        int             `json:"status" bson:"status"`
	VisitTypeList []int           `json:"visit_type_list" bson:"visit_type_list"`
	Marks         []int           `json:"marks" bson:"marks"`
	Keywords      []string        `json:"keywords" bson:"keywords" es:"analyzer:ik_max_word"`
	Areas         []string        `json:"areas" bson:"areas"`
	Sort          int             `json:"sort" bson:"sort"`
	Permission    *PType          `json:"permission,omitempty" bson:"permission,omitempty"`
	Trend         *PostTrend      `json:"trend,omitempty" bson:"trend,omitempty"`
	testAnonymous
	TestMap  map[string]interface{} `json:"test_map,omitempty" bson:"test_map,omitempty"`
	TestMap1 map[string]int         `json:"test_map1,omitempty" bson:"test_map1,omitempty"`
	TestArr  []int                  `json:"test_arr,omitempty" bson:"test_arr,omitempty"`
	TestArr1 []interface{}          `json:"test_arr1,omitempty" bson:"test_arr1,omitempty"`
	TestObj  interface{}            `json:"test_obj,omitempty" bson:"test_obj,omitempty"`
}

func (a *Post) MarshalJSON() ([]byte, error) {
	type Alias Post

	var alias Alias = Alias(*a)

	if alias.Permission != nil {
		for i := 0; i < len(alias.Images); i++ {
			alias.Images[i].Media = nil
		}
		for i := 0; i < len(alias.Videos); i++ {
			alias.Videos[i].Media = nil
		}
		for i := 0; i < len(alias.Audios); i++ {
			alias.Audios[i].Media = nil
		}
	}

	return json.Marshal(alias)
}

type Price struct {
	AssetsType int     `json:"assets_type" bson:"assets_type"`
	Amount     float64 `json:"amount" bson:"amount"`
}

type PType struct {
	Price *Price `json:"price,omitempty" bson:"price,omitempty"`
	Group *GType `json:"group,omitempty" bson:"group,omitempty"`
}

type GType struct {
	Allow *Allow `json:"allow,omitempty" bson:"allow,omitempty"`
	Deny  *Deny  `json:"deny,omitempty" bson:"deny,omitempty"`
}

type Allow struct {
	IDList []string `json:"id_list" bson:"id_list"`
}

type Deny struct {
	IDList []string `json:"id_list" bson:"id_list"`
}
