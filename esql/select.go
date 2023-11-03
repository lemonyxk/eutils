/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:41
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
)

func handleSelect(stmt *sqlparser.Select) (dsl string, table string, err error) {

	var query = M{}

	var result = M{
		"query": query,
	}

	var tableName = String(stmt.From)

	fields, err := FormatSelectExpr(stmt.SelectExprs)
	if err != nil {
		return "", "", err
	}

	if len(fields) != 0 {
		result["_source"] = fields
	}

	// handle limit
	var limit = stmt.Limit
	if limit != nil {
		result["size"] = Number(limit.Rowcount)
		result["from"] = Number(limit.Offset)
	} else {
		result["size"] = 1
		result["from"] = 0
	}

	// handle order by
	var orders = A{}
	for i := 0; i < len(stmt.OrderBy); i++ {
		var orderBy = stmt.OrderBy[i]
		var key = String(orderBy.Expr)
		var val = "asc"
		if orderBy.Direction == sqlparser.DescScr {
			val = "desc"
		}
		orders = append(orders, M{
			key: val,
		})
		result["sort"] = orders
	}

	// handle where
	var where = stmt.Where
	if where == nil {
		delete(result, "query")
		return result.String(), tableName, nil
	}
	handleWhere(query, where.Expr)

	return result.String(), tableName, nil
}

func handleWhere(result M, expr sqlparser.Expr) {
	switch expr.(type) {
	case *sqlparser.AndExpr:
		var query = &A{}
		result["bool"] = M{"filter": query}
		handleExpr(result, query, expr.(*sqlparser.AndExpr).Left, expr)
		handleExpr(result, query, expr.(*sqlparser.AndExpr).Right, expr)
	case *sqlparser.OrExpr:
		var query = &A{}
		result["bool"] = M{"should": query}
		handleExpr(result, query, expr.(*sqlparser.OrExpr).Left, expr)
		handleExpr(result, query, expr.(*sqlparser.OrExpr).Right, expr)
	case *sqlparser.ParenExpr:
		handleWhere(result, expr.(*sqlparser.ParenExpr).Expr)
	default:
		var query = &A{}
		result["bool"] = M{"filter": query}
		handleExpr(result, query, expr, nil)
	}
}

func handleExpr(result M, query *A, expr sqlparser.Expr, parent sqlparser.Expr) {
	switch expr.(type) {
	case *sqlparser.ComparisonExpr:
		handleComparison(query, expr.(*sqlparser.ComparisonExpr))
	case *sqlparser.IsExpr:
		handleIs(query, expr.(*sqlparser.IsExpr))
	case *sqlparser.RangeCond:
		handleRange(query, expr.(*sqlparser.RangeCond))
	case *sqlparser.AndExpr:
		if _, ok := parent.(*sqlparser.AndExpr); ok {
			handleExpr(result, query, expr.(*sqlparser.AndExpr).Left, expr)
			handleExpr(result, query, expr.(*sqlparser.AndExpr).Right, expr)
		} else {
			var res = M{}
			*query = append(*query, res)
			handleAnd(res, expr.(*sqlparser.AndExpr))
		}
	case *sqlparser.OrExpr:
		if _, ok := parent.(*sqlparser.OrExpr); ok {
			handleExpr(result, query, expr.(*sqlparser.OrExpr).Left, expr)
			handleExpr(result, query, expr.(*sqlparser.OrExpr).Right, expr)
		} else {
			var res = M{}
			*query = append(*query, res)
			handleOr(res, expr.(*sqlparser.OrExpr))
		}
	case *sqlparser.ParenExpr:
		handleExpr(result, query, expr.(*sqlparser.ParenExpr).Expr, parent)
	case *sqlparser.FuncExpr:
		handleFunc(result, expr.(*sqlparser.FuncExpr))
	default:
		panic("not support " + String(expr))
	}
}

func handleRange(result *A, cond *sqlparser.RangeCond) {
	var field = sqlparser.String(cond.Left)
	var from = FormatSingle(cond.From)
	var to = FormatSingle(cond.To)
	var query = M{
		"bool": M{
			"filter": M{
				"range": M{
					field: M{
						"gte": from,
						"lte": to,
					},
				},
			},
		},
	}
	*result = append(*result, query)
}

func handleAnd(result M, expr *sqlparser.AndExpr) {
	var query = &A{}
	result["bool"] = M{"filter": query}
	handleExpr(result, query, expr.Left, expr)
	handleExpr(result, query, expr.Right, expr)
}

func handleOr(result M, expr *sqlparser.OrExpr) {
	var query = &A{}
	result["bool"] = M{"should": query}
	handleExpr(result, query, expr.Left, expr)
	handleExpr(result, query, expr.Right, expr)
}

func handleFunc(result M, expr *sqlparser.FuncExpr) {
	var query = &A{}
	result["bool"].(M)["should"] = query

	stmt, err := sqlparser.Parse("select x from a where " + String(expr.Exprs[0]))
	if err != nil {
		panic(err)
	}
	handleExpr(result, query, stmt.(*sqlparser.Select).Where.Expr, expr)
}
