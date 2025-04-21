/**
* @program: engine
*
* @create: 2024-08-30 17:24
**/

package elastic

import "github.com/lemonyxk/kitty/kitty"

type QueryTable struct {
	Name string `json:"name,omitempty"`
	Query
}

type Query struct {
	Skip

	Pre []Q `json:"pre,omitempty"`
	And []Q `json:"and,omitempty"`
	Or  []Q `json:"or,omitempty"`

	Should []Q    `json:"should,omitempty"`
	Sort   []Sort `json:"sort,omitempty"`

	SQL SQL `json:"sql,omitempty"`

	DSL map[string]any `json:"dsl,omitempty"`

	Includes []string `json:"includes,omitempty"`
	Excludes []string `json:"excludes,omitempty"`
	Indexes  []string `json:"indexes,omitempty"`
}

type Skip struct {
	Page  int `json:"page,omitempty"`
	Limit int `json:"limit,omitempty"`
}

type Sort struct {
	Field string `json:"field,omitempty"`
	Order string `json:"order,omitempty"`
}

type Q struct {
	Field string `json:"field,omitempty"`
	Op    string `json:"op,omitempty"`
	Value any    `json:"value,omitempty"`
	Boost int    `json:"boost,omitempty"`
}

type Script struct {
	ID     string `json:"id,omitempty"`
	Source string `json:"source,omitempty"`
	Lang   string `json:"lang,omitempty"`
	Params any    `json:"params,omitempty"`
}

type Unset struct {
	Target string `json:"target,omitempty"`
}

type Replace struct {
	Target string `json:"target,omitempty"`
	Old    any    `json:"old,omitempty"`
	New    any    `json:"new,omitempty"`
}

type Pull struct {
	Target string `json:"target,omitempty"`
	Value  any    `json:"value,omitempty"`
}

type Add struct {
	Target string `json:"target,omitempty"`
	Value  any    `json:"value,omitempty"`
}

type Remove struct {
	Target string `json:"target,omitempty"`
	Index  int    `json:"index,omitempty"`
}

type Params struct {
	Set     kitty.M  `json:"$set,omitempty" bson:"$set,omitempty"`
	Inc     kitty.M  `json:"$inc,omitempty" bson:"$inc,omitempty"`
	Unset   *Unset   `json:"$unset,omitempty" bson:"-"`
	Replace *Replace `json:"$replace,omitempty" bson:"-"`
	Pull    *Pull    `json:"$pull,omitempty" bson:"-"`
	Add     *Add     `json:"$add,omitempty" bson:"-"`
	Remove  *Remove  `json:"$remove,omitempty" bson:"-"`
}
