/**
* @program: es
*
* @description:
*
* @author: lemo
*
* @create: 2023-05-21 18:17
**/

package main

type Media struct {
	Resource Resource `json:"resource" bson:"resource" mapstructure:"resource" es:"type:nested"`
	File     File     `json:"file" bson:"file" mapstructure:"file"`
	Link     Link     `json:"link" bson:"link" mapstructure:"link"`
}

type MediaWithSize struct {
	Media `bson:"inline"`
	Size  Size `json:"size" bson:"size" mapstructure:"size"`
}

type Size struct {
	Width  int `json:"width" bson:"width" mapstructure:"width"`
	Height int `json:"height" bson:"height" mapstructure:"height"`
}

type Link struct {
	Domain string `json:"domain" bson:"domain" mapstructure:"domain"`
	Params string `json:"params" bson:"params" mapstructure:"params"`
}

type Resource struct {
	Path        string `json:"path" bson:"path" mapstructure:"path"`
	Name        string `json:"name" bson:"name" mapstructure:"name"`
	StorageType int    `json:"storage_type" bson:"storage_type" mapstructure:"storage_type"`
}

type File struct {
	Name string `json:"name" bson:"name" mapstructure:"name"`
	Size int64  `json:"size" bson:"size" mapstructure:"size"`
	Mime string `json:"mime" bson:"mime" mapstructure:"mime"`
}

type PostRecord struct {
	Like    int `json:"like" bson:"like" mapstructure:"like"`
	UnLike  int `json:"unlike" bson:"unlike" mapstructure:"unlike"`
	Comment int `json:"comment" bson:"comment" mapstructure:"comment"`
	Share   int `json:"share" bson:"share" mapstructure:"share"`
	Read    int `json:"read" bson:"read" mapstructure:"read"`
	Collect int `json:"collect" bson:"collect" mapstructure:"collect"`
}

type User struct {
	ID        string `json:"id" bson:"_id" mapstructure:"id"`
	UUID      string `json:"uuid" bson:"uuid" mapstructure:"uuid" index:"uuid_1,unique"`
	Pid       string `json:"pid" bson:"pid" mapstructure:"pid" index:"pid_1"`
	UserName  string `json:"user_name" bson:"user_name" mapstructure:"user_name" index:"user_name_1,unique"`
	PackageID int    `json:"package_id" bson:"package_id" mapstructure:"package_id" index:"package_id_1"`

	LoginTime int64  `json:"login_time" bson:"login_time" mapstructure:"login_time" index:"login_time_1"`
	LoginIP   string `json:"login_ip" bson:"login_ip" mapstructure:"login_ip"`
	ReginTime int64  `json:"regin_time" bson:"regin_time" mapstructure:"regin_time" index:"regin_time_1"`
	ReginIP   string `json:"regin_ip" bson:"regin_ip" mapstructure:"regin_ip"`

	Status int `json:"status" bson:"status" mapstructure:"status"`
}

type Post struct {
	ID           string          `json:"id" bson:"_id" mapstructure:"id" es:"type:keyword"`
	Index        string          `json:"index" bson:"index" mapstructure:"index" index:"index_1,unique"`
	UserID       string          `json:"user_id" bson:"user_id" mapstructure:"user_id" index:"user_id_1"`
	User         *User           `json:"user,omitempty" bson:"-" mapstructure:"user" es:"user"`
	Type         int             `json:"type" bson:"type" mapstructure:"type" index:"type_1"`
	CategoriesID []string        `json:"categories_id" bson:"categories_id" mapstructure:"categories_id" index:"categories_id_1"`
	Tags         []string        `json:"tags" bson:"tags" mapstructure:"tags" index:"tags_1"`
	Title        string          `json:"title" bson:"title" mapstructure:"title" index:"title_1"`
	Desc         string          `json:"desc" bson:"desc" mapstructure:"desc"`
	Content      string          `json:"content" bson:"content" mapstructure:"content"`
	Covers       []MediaWithSize `json:"covers" bson:"covers" mapstructure:"covers"`
	Videos       []Media         `json:"videos" bson:"videos" mapstructure:"videos"`
	Images       []Media         `json:"images" bson:"images" mapstructure:"images"`
	Audios       []Media         `json:"audios" bson:"audios" mapstructure:"audios" es:"type:keyword"`
	CreateTime   int64           `json:"create_time" bson:"create_time" mapstructure:"create_time" index:"create_time_1"`
	UpdateTime   int64           `json:"update_time" bson:"update_time" mapstructure:"update_time" index:"update_time_1"`
	Status       int             `json:"status" bson:"status" mapstructure:"status"`

	Record PostRecord `json:"record" bson:"record" mapstructure:"record"`

	Map map[string]interface{} `json:"map" bson:"map" mapstructure:"map" es:"type:keyword"`

	Array [5]int `json:"array" bson:"array" mapstructure:"array"`

	Slice []interface{} `json:"slice" bson:"slice" mapstructure:"slice" es:"type:keyword"`
}

func (a *Post) Empty() bool {
	return a.ID == ""
}
