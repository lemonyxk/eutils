/**
* @program: eutils
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-23 18:05
**/

package main

type Account struct {
	ID         string    `json:"id" bson:"_id"`
	UUID       string    `json:"uuid" bson:"uuid" index:"uuid_1,unique"`
	PassWord   string    `json:"pass_word" bson:"pass_word"`
	PackageID  int       `json:"package_id" bson:"package_id"`
	Wallets    []Wallet  `json:"wallets,omitempty" bson:"wallets,omitempty"`
	Property   *Property `json:"property,omitempty" bson:"property,omitempty"`
	Assets     *Assets   `json:"assets,omitempty" bson:"assets,omitempty"`
	Fraction   *Fraction `json:"fraction,omitempty" bson:"fraction,omitempty"`
	Profile    *Profile  `json:"profile,omitempty" bson:"profile,omitempty"`
	GeoIP      *GeoIP    `json:"geo_ip,omitempty" bson:"geo_ip,omitempty"`
	Info       *Info     `json:"info,omitempty" bson:"info,omitempty" index:"info.email_1"`
	Relation   *Relation `json:"relation,omitempty" bson:"relation,omitempty"`
	Status     int       `json:"status" bson:"status"`
	CreateTime int64     `json:"create_time" bson:"create_time" index:"create_time_1"`
}

type Assets struct {
	Gold       float64 `json:"gold" bson:"gold"`
	Balance    float64 `json:"balance" bson:"balance"`
	CNY        float64 `json:"cny" bson:"cny"`
	Winnings   float64 `json:"winnings" bson:"winnings"`
	Contribute float64 `json:"contribute" bson:"contribute"`
	Prestige   float64 `json:"prestige" bson:"prestige"`
	Develop    float64 `json:"develop" bson:"develop"`
	Silver     float64 `json:"silver" bson:"silver"`
	Diamond    float64 `json:"diamond" bson:"diamond"`
}

type Wallet struct {
	ID            string `json:"id" bson:"id"`
	Name          string `json:"name" bson:"name"`
	FinancialType int    `json:"financial_type" bson:"financial_type"`
	AuthName      string `json:"auth_name" bson:"auth_name"`
	Address       string `json:"address" bson:"address"`
	Network       string `json:"network" bson:"network"`
}

type Property struct {
	Manager *Manager `json:"manager,omitempty" bson:"manager,omitempty"` // 管理组
	User    *User    `json:"user,omitempty" bson:"user,omitempty"`       // 用户组
	Member  *Member  `json:"member,omitempty" bson:"member,omitempty"`   // 会员组
	Extends *Extend  `json:"extends,omitempty" bson:"extends,omitempty"`
}

type Code int

type RuleType int

type Values []int

type KVS []string

const (
	Bool     RuleType = iota + 1 // values [0:false 1:true]
	Range                        // values [min,max,value]
	Duration                     // values [duration min,max,value value min,max,value]
	Enum                         // values [1,2,3,4,5]
)

type Type struct {
	Code   Code     `json:"code" bson:"code"`
	Type   RuleType `json:"type" bson:"type"`
	Values Values   `json:"values,omitempty" bson:"values,omitempty"`
	Enum   KVS      `json:"enum,omitempty" bson:"enum,omitempty"`
}

type Object struct {
	*Type `json:",omitempty" bson:"inline,omitempty"`
	// MAX_INT32 = 2147483647
	ExpireTime int64 `json:"expire_time,omitempty" bson:"expire_time,omitempty"`
}

type Objects map[Code]*Object

type Manager struct {
	Objects    Objects `json:"objects" bson:"objects"`
	ExpireTime int64   `json:"expire_time,omitempty" bson:"expire_time,omitempty"`
}

type User struct {
	Objects    Objects `json:"objects" bson:"objects"`
	ExpireTime int64   `json:"expire_time,omitempty" bson:"expire_time,omitempty"`
}

type Member struct {
	Objects    Objects `json:"objects" bson:"objects"`
	ExpireTime int64   `json:"expire_time,omitempty" bson:"expire_time,omitempty"`
}

type Extend struct {
	Objects `json:",omitempty" bson:"inline,omitempty"`
}

type Fraction struct {
	Level int   `json:"level" bson:"level"`
	Score int64 `json:"score" bson:"score"`
	Need  int64 `json:"need" bson:"need"`
	Count int64 `json:"count" bson:"count"`
}

type Profile struct {
	Bio      string `json:"bio" bson:"bio"`
	Address  string `json:"address" bson:"address"`
	Gender   int    `json:"gender" bson:"gender"`
	Birthday int64  `json:"birthday" bson:"birthday"`
}

type GeoIP struct {
	Ip       string `json:"ip"`
	Hostname string `json:"hostname"`
	City     string `json:"city"`
	Region   string `json:"region"`
	Country  string `json:"country"`
	Loc      string `json:"loc"`
	Org      string `json:"org"`
	Postal   string `json:"postal"`
	Timezone string `json:"timezone"`
}

type Relation struct {
	Blacklist []string `json:"blacklist" bson:"blacklist"`
}
