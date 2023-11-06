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
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)
		val = strings.ReplaceAll(str, "%", "")
		v["query"] = val
		*result = append(*result, M{
			"bool": M{
				"must": A{
					M{
						"match_phrase": M{
							left: v,
						},
					},
				},
			},
		})
	case sqlparser.NotLikeStr: // not like
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)
		val = strings.ReplaceAll(str, "%", "")
		v["query"] = val
		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						"match_phrase": M{
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
