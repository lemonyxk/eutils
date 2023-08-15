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

import "github.com/xwb1989/sqlparser"

func handleIs(result *A, expr *sqlparser.IsExpr) {
	switch expr.Operator {
	case sqlparser.IsNullStr:
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"exists": M{
						"field": String(expr.Expr),
					},
				},
			},
		})
	case sqlparser.IsNotNullStr:
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"exists": M{
						"field": String(expr.Expr),
					},
				},
			},
		})
	case sqlparser.IsTrueStr:
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"term": M{
						String(expr.Expr): true,
					},
				},
			},
		})
	case sqlparser.IsFalseStr:
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"term": M{
						String(expr.Expr): false,
					},
				},
			},
		})
	case sqlparser.IsNotTrueStr:
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"term": M{
						String(expr.Expr): true,
					},
				},
			},
		})
	case sqlparser.IsNotFalseStr:
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"term": M{
						String(expr.Expr): false,
					},
				},
			},
		})
	}
}
