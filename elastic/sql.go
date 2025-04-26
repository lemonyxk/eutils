package elastic

import (
	"fmt"
	"github.com/lemonyxk/eutils/esql"
	"github.com/lemonyxk/kitty/json"
	"reflect"
	"strings"

	"errors"
)

type SQL string

var (
	maps = map[string]string{
		"$eq":      "=",
		"$ne":      "!=",
		"$gt":      ">",
		"$gte":     ">=",
		"$lt":      "<",
		"$lte":     "<=",
		"$in":      "IN",
		"$nin":     "NOT IN",
		"$like":    "LIKE",
		"$nlike":   "NOT LIKE",
		"$between": "BETWEEN",
		"$exists":  "IS NOT NULL",
		"$nex":     "IS NULL",
	}
)

func fix(s string) string {
	if v, ok := maps[s]; ok {
		return v
	}
	return s
}

func MakeSQL(queryRequest QueryTable) (SQL, error) {
	// if len(queryRequest.Params) == 0 {
	// 	return "", errors.New("query is empty")
	// }

	if queryRequest.Page == 0 {
		queryRequest.Page = 1
	}

	if queryRequest.Limit < 0 {
		return "", errors.New("limit must be greater than 0")
	}

	// if queryRequest.Page > 100 {
	// 	return "", errors.New("page is too large")
	// }

	if queryRequest.Limit > 10000 {
		return "", errors.New("limit is too large")
	}

	if queryRequest.Limit*(queryRequest.Page-1) > 10000 {
		return "", errors.New("page is too large")
	}

	if len(queryRequest.Should) > 0 {
		queryRequest.Sort = append([]Sort{{"_score", "DESC"}}, queryRequest.Sort...)
	}

	if len(queryRequest.DSL) > 0 {
		queryRequest.DSL["size"] = queryRequest.Limit
		queryRequest.DSL["from"] = (queryRequest.Page - 1) * queryRequest.Limit
		var bts, err = json.Marshal(queryRequest.DSL)
		return SQL(bts), err
	}

	var sql = fmt.Sprintf("SELECT * FROM `%s`", queryRequest.Name)

	if len(queryRequest.Includes) > 0 {
		sql = fmt.Sprintf("SELECT `includes(%s)` FROM `%s`", strings.Join(queryRequest.Includes, ", "), queryRequest.Name)
	}

	if len(queryRequest.Excludes) > 0 {
		sql = fmt.Sprintf("SELECT `excludes(%s)` FROM `%s`", strings.Join(queryRequest.Excludes, ", "), queryRequest.Name)
	}

	var sub string
	var pre []string
	var and []string
	var or []string
	var should []string

	for j := 0; j < len(queryRequest.Pre); j++ {
		var query = queryRequest.Pre[j]
		if query.Empty() {
			return "", errors.New("pre is empty")
		}

		query.Op = strings.ToUpper(fix(query.Op))

		if j > 0 {
			pre = append(pre, "AND")
		}

		if query.Op == "BETWEEN" {
			var t = reflect.TypeOf(query.Value)
			var v = reflect.ValueOf(query.Value)
			var res []any
			if t.Kind() != reflect.Slice {
				return "", errors.New("between value must be a slice")
			}
			if v.Len() != 2 {
				return "", errors.New("between value length must be 2")
			}
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				res = append(res, fmt.Sprintf("%v", n))
			}
			if query.Boost > 0 {
				pre = append(pre, fmt.Sprintf("`%s#%d` BETWEEN %v AND %v", query.Field, query.Boost, res[0], res[1]))
			} else {
				pre = append(pre, fmt.Sprintf("`%s` BETWEEN %v AND %v", query.Field, res[0], res[1]))
			}
			continue
		}

		if query.Op == "IS NOT NULL" || query.Op == "IS NULL" {
			if query.Boost > 0 {
				pre = append(pre, fmt.Sprintf("`%s#%d` %s", query.Field, query.Boost, query.Op))
			} else {
				pre = append(pre, fmt.Sprintf("`%s` %s", query.Field, query.Op))
			}
			continue
		}

		var t = reflect.TypeOf(query.Value)
		switch t.Kind() {
		case reflect.Slice:
			var in = ""
			var v = reflect.ValueOf(query.Value)
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				if k == 0 {
					in += fmt.Sprintf("%v", n)
				} else {
					in += fmt.Sprintf(",%v", n)
				}
			}
			if query.Boost > 0 {
				pre = append(pre, fmt.Sprintf("`%s#%d` %s (%v)", query.Field, query.Boost, query.Op, in))
			} else {
				pre = append(pre, fmt.Sprintf("`%s` %s (%v)", query.Field, query.Op, in))
			}
		default:
			var v = query.Value
			var dt = reflect.TypeOf(query.Value)
			if dt.Kind() == reflect.String {
				v = fmt.Sprintf("'%v'", v)
			}
			if dt.Kind() == reflect.Float64 {
				// is integer
				if v.(float64) == float64(int(v.(float64))) {
					v = int(v.(float64))
				}
			}
			if query.Boost > 0 {
				pre = append(pre, fmt.Sprintf("`%s#%d` %s %v", query.Field, query.Boost, query.Op, v))
			} else {
				pre = append(pre, fmt.Sprintf("`%s` %s %v", query.Field, query.Op, v))
			}
		}
	}

	for j := 0; j < len(queryRequest.And); j++ {
		var query = queryRequest.And[j]
		if query.Empty() {
			return "", errors.New("and is empty")
		}

		query.Op = strings.ToUpper(fix(query.Op))

		if j > 0 {
			and = append(and, "AND")
		}

		if query.Op == "BETWEEN" {
			var t = reflect.TypeOf(query.Value)
			var v = reflect.ValueOf(query.Value)
			var res []any
			if t.Kind() != reflect.Slice {
				return "", errors.New("between value must be a slice")
			}
			if v.Len() != 2 {
				return "", errors.New("between value length must be 2")
			}
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				res = append(res, fmt.Sprintf("%v", n))
			}
			if query.Boost > 0 {
				and = append(and, fmt.Sprintf("`%s#%d` BETWEEN %v AND %v", query.Field, query.Boost, res[0], res[1]))
			} else {
				and = append(and, fmt.Sprintf("`%s` BETWEEN %v AND %v", query.Field, res[0], res[1]))
			}
			continue
		}

		if query.Op == "IS NOT NULL" || query.Op == "IS NULL" {
			if query.Boost > 0 {
				and = append(and, fmt.Sprintf("`%s#%d` %s", query.Field, query.Boost, query.Op))
			} else {
				and = append(and, fmt.Sprintf("`%s` %s", query.Field, query.Op))
			}
			continue
		}

		var t = reflect.TypeOf(query.Value)
		switch t.Kind() {
		case reflect.Slice:
			var in = ""
			var v = reflect.ValueOf(query.Value)
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				if k == 0 {
					in += fmt.Sprintf("%v", n)
				} else {
					in += fmt.Sprintf(",%v", n)
				}
			}
			if query.Boost > 0 {
				and = append(and, fmt.Sprintf("`%s#%d` %s (%v)", query.Field, query.Boost, query.Op, in))
			} else {
				and = append(and, fmt.Sprintf("`%s` %s (%v)", query.Field, query.Op, in))
			}
		default:
			var v = query.Value
			var dt = reflect.TypeOf(query.Value)
			if dt.Kind() == reflect.String {
				v = fmt.Sprintf("'%v'", v)
			}
			if dt.Kind() == reflect.Float64 {
				// is integer
				if v.(float64) == float64(int(v.(float64))) {
					v = int(v.(float64))
				}
			}
			if query.Boost > 0 {
				and = append(and, fmt.Sprintf("`%s#%d` %s %v", query.Field, query.Boost, query.Op, v))
			} else {
				and = append(and, fmt.Sprintf("`%s` %s %v", query.Field, query.Op, v))
			}
		}
	}

	for j := 0; j < len(queryRequest.Or); j++ {
		var query = queryRequest.Or[j]
		if query.Empty() {
			return "", errors.New("or is empty")
		}

		query.Op = strings.ToUpper(fix(query.Op))

		if j > 0 {
			or = append(or, "OR")
		}

		if query.Op == "BETWEEN" {
			var t = reflect.TypeOf(query.Value)
			var v = reflect.ValueOf(query.Value)
			var res []any
			if t.Kind() != reflect.Slice {
				return "", errors.New("between value must be a slice")
			}
			if v.Len() != 2 {
				return "", errors.New("between value length must be 2")
			}
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				res = append(res, fmt.Sprintf("%v", n))
			}
			if query.Boost > 0 {
				or = append(or, fmt.Sprintf("`%s#%d` BETWEEN %v AND %v", query.Field, query.Boost, res[0], res[1]))
			} else {
				or = append(or, fmt.Sprintf("`%s` BETWEEN %v AND %v", query.Field, res[0], res[1]))
			}
			continue
		}

		if query.Op == "IS NOT NULL" || query.Op == "IS NULL" {
			if query.Boost > 0 {
				or = append(or, fmt.Sprintf("`%s#%d` %s", query.Field, query.Boost, query.Op))
			} else {
				or = append(or, fmt.Sprintf("`%s` %s", query.Field, query.Op))
			}
			continue
		}

		var t = reflect.TypeOf(query.Value)
		switch t.Kind() {
		case reflect.Slice:
			var in = ""
			var v = reflect.ValueOf(query.Value)
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				if k == 0 {
					in += fmt.Sprintf("%v", n)
				} else {
					in += fmt.Sprintf(",%v", n)
				}
			}
			if query.Boost > 0 {
				or = append(or, fmt.Sprintf("`%s#%d` %s (%v)", query.Field, query.Boost, query.Op, in))
			} else {
				or = append(or, fmt.Sprintf("`%s` %s (%v)", query.Field, query.Op, in))
			}
		default:
			var v = query.Value
			var dt = reflect.TypeOf(query.Value)
			if dt.Kind() == reflect.String {
				v = fmt.Sprintf("'%v'", v)
			}
			if dt.Kind() == reflect.Float64 {
				// is integer
				if v.(float64) == float64(int(v.(float64))) {
					v = int(v.(float64))
				}
			}
			if query.Boost > 0 {
				or = append(or, fmt.Sprintf("`%s#%d` %s %v", query.Field, query.Boost, query.Op, v))
			} else {
				or = append(or, fmt.Sprintf("`%s` %s %v", query.Field, query.Op, v))
			}
		}
	}

	for j := 0; j < len(queryRequest.Should); j++ {
		var query = queryRequest.Should[j]
		if query.Empty() {
			return "", errors.New("should is empty")
		}

		query.Op = strings.ToUpper(fix(query.Op))

		if j > 0 {

		}

		if query.Op == "BETWEEN" {
			var t = reflect.TypeOf(query.Value)
			var v = reflect.ValueOf(query.Value)
			var res []any
			if t.Kind() != reflect.Slice {
				return "", errors.New("between value must be a slice")
			}
			if v.Len() != 2 {
				return "", errors.New("between value length must be 2")
			}
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if n == nil {
					return "", errors.New("slice value must not be nil")
				}
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				res = append(res, fmt.Sprintf("%v", n))
			}
			if query.Boost > 0 {
				should = append(should, fmt.Sprintf("`%s#%d` BETWEEN %v AND %v", query.Field, query.Boost, res[0], res[1]))
			} else {
				should = append(should, fmt.Sprintf("`%s` BETWEEN %v AND %v", query.Field, res[0], res[1]))
			}
			continue
		}

		if query.Op == "IS NOT NULL" || query.Op == "IS NULL" {
			if query.Boost > 0 {
				should = append(should, fmt.Sprintf("`%s#%d` %s", query.Field, query.Boost, query.Op))
			} else {
				should = append(should, fmt.Sprintf("`%s` %s", query.Field, query.Op))
			}
			continue
		}

		var t = reflect.TypeOf(query.Value)
		switch t.Kind() {
		case reflect.Slice:
			var in = ""
			var v = reflect.ValueOf(query.Value)
			for k := 0; k < v.Len(); k++ {
				var n = v.Index(k).Interface()
				if reflect.TypeOf(n).Kind() == reflect.String {
					n = fmt.Sprintf("'%v'", n)
				}
				if reflect.TypeOf(n).Kind() == reflect.Float64 {
					// is integer
					if n.(float64) == float64(int(n.(float64))) {
						n = int(n.(float64))
					}
				}
				if k == 0 {
					in += fmt.Sprintf("%v", n)
				} else {
					in += fmt.Sprintf(",%v", n)
				}
			}
			if query.Boost > 0 {
				should = append(should, fmt.Sprintf("`%s#%d` %s (%v)", query.Field, query.Boost, query.Op, in))
			} else {
				should = append(should, fmt.Sprintf("`%s` %s (%v)", query.Field, query.Op, in))
			}
		default:
			var v = query.Value
			var dt = reflect.TypeOf(query.Value)
			if dt.Kind() == reflect.String {
				v = fmt.Sprintf("'%v'", v)
			}
			if dt.Kind() == reflect.Float64 {
				// is integer
				if v.(float64) == float64(int(v.(float64))) {
					v = int(v.(float64))
				}
			}
			if query.Boost > 0 {
				should = append(should, fmt.Sprintf("`%s#%d` %s %v", query.Field, query.Boost, query.Op, v))
			} else {
				should = append(should, fmt.Sprintf("`%s` %s %v", query.Field, query.Op, v))
			}
		}
	}

	for j := 0; j < len(queryRequest.Sort); j++ {
		var sort = queryRequest.Sort[j]

		if j == 0 {
			sub += " ORDER BY"
		} else {
			sub += ","
		}

		sub += fmt.Sprintf(" `%s` %s", sort.Field, strings.ToUpper(sort.Order))
	}

	if len(queryRequest.SearchAfter) > 0 {
		if len(queryRequest.Sort) == 0 {
			sub += " ORDER BY SEARCH_AFTER("
		} else {
			sub += ", SEARCH_AFTER("
		}
		for i := 0; i < len(queryRequest.SearchAfter); i++ {
			var v = reflect.ValueOf(queryRequest.SearchAfter[i])
			var n = v.Interface()
			if reflect.TypeOf(n).Kind() == reflect.String {
				n = fmt.Sprintf("'%v'", n)
			}
			if reflect.TypeOf(n).Kind() == reflect.Float64 {
				// is integer
				if n.(float64) == float64(int(n.(float64))) {
					n = int(n.(float64))
				}
			}
			if i == 0 {
				sub += fmt.Sprintf("%v", n)
			} else {
				sub += fmt.Sprintf(",%v", n)
			}
		}
		sub += ")"
	}

	sub = sub + fmt.Sprintf(" LIMIT %d,%d", (queryRequest.Page-1)*queryRequest.Limit, queryRequest.Limit)

	if queryRequest.SQL != "" {
		return SQL(sql + " WHERE " + string(queryRequest.SQL) + sub), nil
	}

	var middles []string

	if len(should) > 0 {
		middles = append(middles, "SHOULD( "+strings.Join(should, ", ")+" )")
	}

	if len(pre) > 0 {
		middles = append(middles, strings.Join(pre, " "))
	}

	if len(and) > 0 {
		middles = append(middles, strings.Join(and, " "))
	}

	if len(or) > 0 {
		middles = append(middles, fmt.Sprintf("( %s )", strings.Join(or, " ")))
	}

	var middle = strings.Join(middles, " AND ")

	if len(middle) > 0 {
		return SQL(sql + " WHERE " + middle + sub), nil
	}

	return SQL(sql + sub), nil
}

func BuildSQL(sql string, args ...any) (SQL, error) {
	if len(args) == 0 {
		return SQL(sql), nil
	}
	var index = 0
	for i := 0; i < len(sql); i++ {
		if sql[i] == '?' {
			if index >= len(args) {
				return "", errors.New("args is empty")
			}

			var t = reflect.TypeOf(args[index])
			switch t.Kind() {
			case reflect.Slice:
				var in = ""
				var v = reflect.ValueOf(args[index])
				for k := 0; k < v.Len(); k++ {
					var n = v.Index(k).Interface()
					if reflect.TypeOf(n).Kind() == reflect.String {
						n = fmt.Sprintf("'%v'", n)
					}
					if k == 0 {
						in += fmt.Sprintf("%v", n)
					} else {
						in += fmt.Sprintf(",%v", n)
					}
				}
				sql = strings.Replace(sql, "?", in, 1)
			default:
				var v = args[index]
				switch v.(type) {
				case string:
					v = fmt.Sprintf("'%v'", v)
				}
				sql = strings.Replace(sql, "?", fmt.Sprintf("%v", args[index]), 1)
			}
			index++
		}
	}
	return SQL(sql), nil
}

func ConvertSQLToDSL(sql SQL) (string, string, error) {
	var dsl, table, err = esql.Convert(string(sql))
	if err != nil {
		fmt.Println("sql:", sql)
	}
	return dsl, table, err
}
