/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-15 00:14
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

func handleIs(result *A, expr *sqlparser.IsExpr, action string) {

	var left = String(expr.Expr)
	var v = M{}
	var arr = strings.Split(left, "#")
	if len(arr) == 2 {
		left = arr[0]
		v["boost"] = StringToInt(arr[1])
	}

	switch expr.Operator {
	case sqlparser.IsNullStr:
		v["field"] = left
		*result = append(*result, M{
			"bool": M{
				"must_not": A{
					M{
						"exists": v,
					},
				},
			},
		})
	case sqlparser.IsNotNullStr:
		v["field"] = left
		*result = append(*result, M{
			"bool": M{
				action: A{
					M{
						"exists": v,
					},
				},
			},
		})
	case sqlparser.IsTrueStr:
		v["value"] = true
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
	case sqlparser.IsFalseStr:
		v["value"] = false
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
	case sqlparser.IsNotTrueStr:
		v["value"] = true
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
	case sqlparser.IsNotFalseStr:
		v["value"] = false
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
	}
}
