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
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strings"
	"unicode"
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
		if strings.HasPrefix(strings.ToUpper(key), "CALC(") {
			key = key[5 : len(key)-1]

			var source = doParseCalcField(key)

			orders = append(orders, M{
				"_script": M{
					"type": "number",
					"script": M{
						"source": source,
					},
					"order": val,
				},
			})
		} else if strings.HasPrefix(strings.ToUpper(key), "SCRIPT(") {
			key = key[7 : len(key)-1]

			orders = append(orders, M{
				"_script": M{
					"type": "number",
					"script": M{
						"source": key,
					},
					"order": val,
				},
			})

		} else if strings.HasPrefix(strings.ToUpper(key), "SEARCH_AFTER(") {
			key = key[13 : len(key)-1]
			var arr = strings.Split(key, ", ")
			var arr2 = make([]any, 0)
			for j := 0; j < len(arr); j++ {
				if arr[j][0] == '\'' {
					arr2 = append(arr2, arr[j][1:len(arr[j])-1])
				} else {
					arr2 = append(arr2, StringToFloat(arr[j]))
				}
			}
			result["search_after"] = arr2
		} else {
			orders = append(orders, M{
				key: val,
			})
		}

		if len(orders) > 0 {
			result["sort"] = orders
		}
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

func doParseCalcField(key string) string {
	// 解析计算字段
	// 例如：CALC(a+b*c/d * (f-2)) desc
	// 解析为：a+b*c/d * (f-2)
	// 提取字段：a b c d f

	//def trendWatch = doc['trend.watch'].size() == 0 ? 0 : doc['trend.watch'].value;def initTrendWatch = doc['init_trend.watch'].size() == 0 ? 0 : doc['init_trend.watch'].value;return trendWatch + initTrendWatch;

	var n = 0
	var getNameFromAZ = func() string {
		for {
			var repeat = n / 26
			if repeat > 0 {
				var name = string(rune('a' + repeat - 1))
				name += string(rune('a' + n%26))
				n++
				return name
			}
			var name = string(rune('a' + n))
			n++
			return name
		}
	}

	var defArr []string
	var nameArr []string
	//var fields string
	var returnStr = "return "

	var field = ""
	for i := 0; i < len(key); i++ {
		if unicode.IsSpace(rune(key[i])) {
			continue
		}
		//fields += string(key[i])
		if key[i] == '(' || key[i] == ')' || key[i] == '+' || key[i] == '-' || key[i] == '*' || key[i] == '/' {
			if field != "" {
				if !IsNumber(field) {
					//key = strings.ReplaceAll(key, fields[j], "doc['"+fields[j]+"'].value")
					//fields = fields[:len(fields)-len(field)-1]
					//fields += "doc['" + field + "'].value"
					//fields += string(key[i])
					var name = getNameFromAZ()
					nameArr = append(nameArr, name)
					defArr = append(defArr, fmt.Sprintf("def %s = doc['%s'].size() == 0 ? 0 : doc['%s'].value;", name, field, field))
					returnStr += name
					returnStr += string(key[i])
				}
				field = ""
			}
			continue
		}

		field += string(key[i])

		if i == len(key)-1 {
			if !IsNumber(field) {
				//fields = fields[:len(fields)-len(field)]
				//fields += "doc['" + field + "'].value"
				var name = getNameFromAZ()
				nameArr = append(nameArr, name)
				defArr = append(defArr, fmt.Sprintf("def %s = doc['%s'].size() == 0 ? 0 : doc['%s'].value;", name, field, field))
				returnStr += name
			}
		}
	}

	var res = strings.Join(defArr, "")
	res += returnStr + ";"
	//log.Println(res)

	//return fields
	return res
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
		handleComparison(query, expr.(*sqlparser.ComparisonExpr), "filter")
	case *sqlparser.IsExpr:
		handleIs(query, expr.(*sqlparser.IsExpr), "filter")
	case *sqlparser.RangeCond:
		handleRange(query, expr.(*sqlparser.RangeCond), "filter")
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

	var name = expr.Name.String()

	if strings.ToUpper(name) != "SHOULD" {
		panic("not support " + name + " function")
	}

	var query = &A{}
	result["bool"].(M)["should"] = query

	for i := 0; i < len(expr.Exprs); i++ {
		var expr = expr.Exprs[i].(*sqlparser.AliasedExpr).Expr

		switch expr.(type) {
		case *sqlparser.ComparisonExpr:
			handleComparison(query, expr.(*sqlparser.ComparisonExpr), "must")
		case *sqlparser.IsExpr:
			handleIs(query, expr.(*sqlparser.IsExpr), "must")
		case *sqlparser.RangeCond:
			handleRange(query, expr.(*sqlparser.RangeCond), "must")
		default:
			panic("not support " + String(expr))
		}
	}
}
