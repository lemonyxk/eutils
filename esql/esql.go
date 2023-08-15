/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:37
**/

package esql

import (
	"github.com/xwb1989/sqlparser"
)

func Convert(sql string) (dsl string, table string, err error) {
	stmt, err := sqlparser.Parse(sql)

	if err != nil {
		return "", "", err
	}

	//sql valid, start to handle
	switch stmt.(type) {
	case *sqlparser.Select:
		dsl, table, err = handleSelect(stmt.(*sqlparser.Select))
	case *sqlparser.Update:
		return handleUpdate(stmt.(*sqlparser.Update))
	case *sqlparser.Insert:
		return handleInsert(stmt.(*sqlparser.Insert))
	case *sqlparser.Delete:
		return handleDelete(stmt.(*sqlparser.Delete))
	}

	if err != nil {
		return "", "", err
	}

	return dsl, table, nil
}
