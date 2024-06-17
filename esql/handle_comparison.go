/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-15 00:15
**/

package esql

import (
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strings"
)

func handleComparison(result *A, expr *sqlparser.ComparisonExpr, action string) {
	var val = FormatSingle(expr.Right)
	var left = String(expr.Left)
	var v = M{}
	var arr = strings.Split(left, "#")
	if len(arr) == 2 {
		left = arr[0]
		v["boost"] = StringToInt(arr[1])
	}

	switch expr.Operator {
	case sqlparser.EqualStr: // =
		v["value"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"term": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.LessThanStr: // <
		v["lt"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"range": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.LessEqualStr: // <=
		v["lte"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"range": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.GreaterThanStr: // >
		v["gt"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"range": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.GreaterEqualStr: // >=
		v["gte"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"range": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.NotEqualStr: // !=
		v["value"] = val
		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						"term": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.InStr: // in
		var val = FormatMulti(expr.Right)
		v[left] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"terms": v,
					},
				},
			},
		})
	case sqlparser.NotInStr: // not in
		var val = FormatMulti(expr.Right)
		v[left] = val
		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						"terms": v,
					},
				},
			},
		})
	case sqlparser.LikeStr: // like
		// match will split the keyword and match each of them in any order
		// match_phrase will split the keyword and match them in order but must match all of words
		// match_phrase_prefix will split the keyword and match them in order and support prefix in the last word
		// match will split the keyword and match each of them in any order and support prefix or suffix in each word
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)

		var mode string
		var count = strings.Count(str, "%")
		if count == 0 {
			mode = "match_phrase"
		} else if count == 1 {
			if strings.HasPrefix(str, "%") {
				mode = "match_phrase_prefix"
			} else if strings.HasSuffix(str, "%") {
				mode = "match_phrase_prefix"
			} else {
				mode = "wildcard"
			}
		} else {
			if !(str[0] == '%' && str[len(str)-1] == '%') {
				mode = "wildcard"
			}
		}

		if strings.Count(str, "*") == 2 && strings.HasPrefix(str, "*") && strings.HasSuffix(str, "*") {
			mode = "match"
		}

		val = strings.ReplaceAll(str, "%", "")
		val = strings.ReplaceAll(str, "*", "")

		if mode == "wildcard" {
			v["value"] = val
		} else {
			v["query"] = val
		}

		*result = append(*result, M{
			"bool": M{
				"must": A{
					M{
						mode: M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.NotLikeStr: // not like
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)

		var mode string
		var count = strings.Count(str, "%")
		if count == 0 {
			mode = "match_phrase"
		} else if count == 1 {
			if strings.HasPrefix(str, "%") {
				mode = "match_phrase_prefix"
			} else if strings.HasSuffix(str, "%") {
				mode = "match_phrase_prefix"
			} else {
				mode = "wildcard"
			}
		} else {
			if !(str[0] == '%' && str[len(str)-1] == '%') {
				mode = "wildcard"
			}
		}

		if strings.Count(str, "*") == 2 && strings.HasPrefix(str, "*") && strings.HasSuffix(str, "*") {
			mode = "match"
		}

		val = strings.ReplaceAll(str, "%", "")
		val = strings.ReplaceAll(str, "*", "")

		if mode == "wildcard" {
			v["value"] = val
		} else {
			v["query"] = val
		}

		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						mode: M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.RegexpStr: // regexp
		v["value"] = val
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"regexp": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.NotRegexpStr: // not regexp
		v["value"] = val
		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						"regexp": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.JSONExtractOp: // json_extract
		panic("not support json_extract")
	case sqlparser.JSONUnquoteExtractOp: // json_unquote_extract
		panic("not support json_unquote_extract")
	case sqlparser.NullSafeEqualStr: // <=>
		panic("not support <=>")
	}
}
