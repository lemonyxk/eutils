/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:53
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
	"strconv"
	"strings"
)

func Number(expr sqlparser.SQLNode) int {
	var v, _ = strconv.Atoi(sqlparser.String(expr))
	return v
}

func String(expr sqlparser.SQLNode) string {
	var val = sqlparser.String(expr)
	if val[0] == '`' {
		return strings.ReplaceAll(val, "`", "")
	}
	return val
}

func FormatSingle(expr sqlparser.Expr) any {
	var val any
	var rightString = String(expr)
	if rightString[0] == '\'' { // remove quote
		val = rightString[1 : len(rightString)-1]
	} else {
		val = Number(expr)
		if rightString != "0" && val == 0 {
			val = rightString
		}
	}
	return val
}

func FormatMulti(expr sqlparser.Expr) []any {
	var rightString = String(expr)
	rightString = rightString[1 : len(rightString)-1]
	var arr = strings.Split(rightString, ", ")
	var vs = make([]any, 0)
	for i := 0; i < len(arr); i++ {
		if arr[i][0] == '\'' {
			vs = append(vs, arr[i][1:len(arr[i])-1])
		} else {
			var v, _ = strconv.Atoi(arr[i])
			vs = append(vs, v)
		}
	}
	return vs
}
