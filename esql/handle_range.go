/**
* @program: eutils
*
* @description:
*
* @author: lemo
*
* @create: 2023-11-03 18:54
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
	"strings"
)

func handleRange(result *A, cond *sqlparser.RangeCond, action string) {

	var left = String(cond.Left)
	var v = M{}
	var arr = strings.Split(left, "#")
	if len(arr) == 2 {
		left = arr[0]
		v["boost"] = StringToInt(arr[1])
	}

	var from = FormatSingle(cond.From)
	var to = FormatSingle(cond.To)

	v["gte"] = from
	v["lte"] = to

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
}
