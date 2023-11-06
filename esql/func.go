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
	"errors"
	"fmt"
	"github.com/xwb1989/sqlparser"
	"strconv"
	"strings"
)

func StringToInt(str string) int {
	var v, _ = strconv.Atoi(str)
	return v
}

func Number(expr sqlparser.SQLNode) int {
	var v, _ = strconv.Atoi(sqlparser.String(expr))
	return v
}

func String(expr sqlparser.SQLNode) string {
	var val = sqlparser.String(expr)
	return strings.ReplaceAll(val, "`", "")
}

func FormatSelectExpr(ss sqlparser.SelectExprs) (M, error) {

	var includes []string
	var excludes []string

	for i := 0; i < len(ss); i++ {
		var str = sqlparser.String(ss[i].(sqlparser.SQLNode))
		str = strings.ReplaceAll(str, "`", "")
		if strings.Contains(str, "*") {
			if i != 0 {
				return nil, errors.New("* must be first at select column")
			}
			if str != "*" {
				return nil, errors.New("syntax error at " + str)
			}
			return nil, nil
		}
		if strings.Contains(str, "(") {
			if strings.HasPrefix(strings.ToUpper(str), "INCLUDES(") {
				includes = append(includes, str[10:len(str)-1])
				continue
			}
			if strings.HasPrefix(strings.ToUpper(str), "EXCLUDES(") {
				excludes = append(excludes, str[10:len(str)-1])
				continue
			}
			if str[0] != '(' {
				return nil, fmt.Errorf("do not support func %s", str)
			}
			if strings.Contains(str, ", ") {
				return nil, errors.New("operand should contain 1 column(s)")
			}
			includes = append(includes, str[1:len(str)-1])
			continue
		}
		includes = append(includes, str)
	}

	var result = M{"excludes": excludes, "includes": includes}

	if len(excludes) == 0 {
		delete(result, "excludes")
	}

	if len(includes) == 0 {
		delete(result, "includes")
	}

	return result, nil
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
