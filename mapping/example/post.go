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

type Media struct {
	Resource *Resource `json:"resource,omitempty" bson:"resource,omitempty"`
	File     *File     `json:"file,omitempty" bson:"file,omitempty"`
	Link     *Link     `json:"link,omitempty" bson:"link,omitempty"`
}

type MediaWithItem struct {
	*Media `json:",omitempty" bson:"inline,omitempty"`
	Price  *Price `json:"price,omitempty" bson:"price,omitempty"`
}

type Info struct {
	Duration float64 `json:"duration" bson:"duration"`
}

type Download struct {
	*MediaWithItem `json:",omitempty" bson:"inline,omitempty"`
}

type AudioWithInfo struct {
	*MediaWithItem `json:",omitempty" bson:"inline,omitempty"`
	Info           Info `json:"info" bson:"info"`
}

type VideoWithInfo struct {
	*MediaWithItem `json:",omitempty" bson:"inline,omitempty"`
	Info           Info           `json:"info" bson:"info"`
	Property       *VideoProperty `json:"property,omitempty" bson:"property,omitempty"`
}

type ImageWithSize struct {
	*MediaWithItem `json:",omitempty" bson:"inline,omitempty"`
	Size           Size `json:"size" bson:"size"`
}

type Size struct {
	Width  float64 `json:"width" bson:"width"`
	Height float64 `json:"height" bson:"height"`
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
	Hash string `json:"hash" bson:"hash"`
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

type Line struct {
	Master bool     `json:"master" bson:"master"`
	IDList []string `json:"id_list" bson:"id_list"`
}

type Text struct {
	Content string `json:"content" bson:"content" es:"analyzer:ik_smart"`
	Hash    string `json:"hash" bson:"hash"`
	Price   *Price `json:"price,omitempty" bson:"price,omitempty"`
}

func (a *Text) Empty() bool {
	return a == nil || a.Content == ""
}

type Post struct {
	ID            string          `json:"id" bson:"_id"`
	EID           int64           `json:"eid" bson:"eid" index:"eid_1,unique"`
	UserID        string          `json:"user_id" bson:"user_id" index:"user_id_1"`
	PackageID     int             `json:"package_id" bson:"package_id"`
	Type          int             `json:"type" bson:"type"`
	CategoriesID  []string        `json:"categories_id" bson:"categories_id" index:"categories_id_1"`
	Tags          []string        `json:"tags" bson:"tags" index:"tags_1"`
	Title         string          `json:"title" bson:"title" index:"title_1" es:"analyzer:ik_max_word"`
	Desc          string          `json:"desc" bson:"desc" es:"analyzer:ik_smart"`
	Covers        []ImageWithSize `json:"covers" bson:"covers"`
	CreateTime    int64           `json:"create_time" bson:"create_time" index:"create_time_1"`
	UpdateTime    int64           `json:"update_time" bson:"update_time" index:"update_time_1"`
	Status        int             `json:"status" bson:"status"`
	VisitTypeList []int           `json:"visit_type_list" bson:"visit_type_list"`
	Marks         []int           `json:"marks" bson:"marks"`
	Keywords      []string        `json:"keywords" bson:"keywords" es:"analyzer:ik_max_word"`
	Areas         []string        `json:"areas" bson:"areas"`
	Sort          int             `json:"sort" bson:"sort"`
	Trend         *PostTrend      `json:"trend,omitempty" bson:"trend,omitempty"`
	Relation      *Relation       `json:"relation,omitempty" bson:"relation,omitempty"`
	Value         `bson:"value"`
}

type Value struct {
	Text      []Text          `json:"text,omitempty" bson:"text,omitempty"`
	Images    []ImageWithSize `json:"images,omitempty" bson:"images,omitempty"`
	Videos    []VideoWithInfo `json:"videos,omitempty" bson:"videos,omitempty"`
	Audios    []AudioWithInfo `json:"audios,omitempty" bson:"audios,omitempty"`
	Downloads []Download      `json:"downloads,omitempty" bson:"downloads,omitempty"`
	Price     *Price          `json:"price,omitempty" bson:"price,omitempty"`
}

type Price struct {
	AssetsType int     `json:"assets_type" bson:"assets_type"`
	Amount     float64 `json:"amount" bson:"amount"`
}
