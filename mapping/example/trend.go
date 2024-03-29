/**
* @program: eutils
*
* @description:
*
* @author: lemo
*
* @create: 2023-09-01 16:10
**/

package main

type ID string

type ElasticID int64

type Trend struct {
	ID         ID        `json:"id" bson:"_id"`
	EID        ElasticID `json:"eid" bson:"eid" index:"eid_1,unique"`
	UserID     ID        `json:"user_id" bson:"user_id" index:"user_id_1"`
	ForID      ID        `json:"for_id" bson:"for_id" index:"for_id_1"`
	BelongID   ID        `json:"belong_id" bson:"belong_id" index:"belong_id_1"`
	Type       int       `json:"type" bson:"type"`
	Action     int       `json:"action" bson:"action"`
	Counter    int64     `json:"counter" bson:"counter"`
	ForTime    int64     `json:"for_time" bson:"for_time"`
	CreateTime int64     `json:"create_time" bson:"create_time" index:"create_time_1"`
	Name       float64   `json:"name" bson:"name"`
	Link       float64   `json:"link" bson:"link"`
	Post       `bson:"post"`
	Item
}

type Post1 struct {
	//Name string `json:"name" bson:"name"`
	Link string `json:"link" bson:"link"`
	Post2
}

type Post2 struct {
	Name1 string `json:"name1" bson:"name1"`
	Link1 string `json:"link1" bson:"link1"`
}

type Trend1 struct {
	ID ID `json:"id" bson:"_id"`
	Item1
}

type Item1 struct {
	ID ID `json:"id" bson:"_id"`
	Post1
}
