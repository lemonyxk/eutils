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

func handleComparison(result *A, expr *sqlparser.ComparisonExpr) {
	switch expr.Operator {
	case sqlparser.EqualStr: // =
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"term": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.LessThanStr: // <
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"range": M{
						String(expr.Left): M{
							"lt": val,
						},
					},
				},
			},
		})
	case sqlparser.LessEqualStr: // <=
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"range": M{
						String(expr.Left): M{
							"lte": val,
						},
					},
				},
			},
		})
	case sqlparser.GreaterThanStr: // >
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"range": M{
						String(expr.Left): M{
							"gt": val,
						},
					},
				},
			},
		})
	case sqlparser.GreaterEqualStr: // >=
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"range": M{
						String(expr.Left): M{
							"gte": val,
						},
					},
				},
			},
		})
	case sqlparser.NotEqualStr: // !=
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"term": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.InStr: // in
		var val = FormatMulti(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"terms": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.NotInStr: // not in
		var val = FormatMulti(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"terms": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.LikeStr: // like
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)
		val = strings.ReplaceAll(str, "%", "")
		*result = append(*result, M{
			"match": M{
				String(expr.Left): val,
			},
		})
	case sqlparser.NotLikeStr: // not like
		var val = FormatSingle(expr.Right)
		var str = fmt.Sprintf("%v", val)
		val = strings.ReplaceAll(str, "%", "")
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"match": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.RegexpStr: // regexp
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"filter": M{
					"regexp": M{
						String(expr.Left): val,
					},
				},
			},
		})
	case sqlparser.NotRegexpStr: // not regexp
		var val = FormatSingle(expr.Right)
		*result = append(*result, M{
			"bool": M{
				"must_not": M{
					"regexp": M{
						String(expr.Left): val,
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
