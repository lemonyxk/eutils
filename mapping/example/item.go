/**
* @program: eutils
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-24 11:59
**/

package main

type Context struct {
	MediaList        []Media `json:"media_list,omitempty" bson:"media_list,omitempty"`
	PropertyTypeList []int   `json:"property_type_list,omitempty" bson:"property_type_list,omitempty"`
}

type Item struct {
	ID         string          `json:"id" bson:"_id"`
	PackageID  int             `json:"package_id" bson:"package_id"`
	Name       string          `json:"name" bson:"name" index:"name_1"`
	Type       int             `json:"type" bson:"type"`
	Tags       []string        `json:"tags" bson:"tags" index:"tags_1"`
	Expire     Expire          `json:"expire" bson:"expire"`
	Price      Price           `json:"price" bson:"price" index:"price.amount_1"`
	Images     []ImageWithSize `json:"images" bson:"images"`
	Context    Context         `json:"context" bson:"context"`
	CreateTime int64           `json:"create_time" bson:"create_time" index:"create_time_1"`
	UpdateTime int64           `json:"update_time" bson:"update_time" index:"update_time_1"`
	StartTime  int64           `json:"start_time" bson:"start_time" index:"start_time_1"`
	EndTime    int64           `json:"end_time" bson:"end_time" index:"end_time_1"`
	Sort       int64           `json:"sort" bson:"sort"`
	Status     int             `json:"status" bson:"status"`
}

type Expire struct {
	Value int `json:"value" bson:"value"`
	Code  int `json:"code" bson:"code"`
}
