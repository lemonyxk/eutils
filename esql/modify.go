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
	"errors"
	"github.com/xwb1989/sqlparser"
)

func handleUpdate(update *sqlparser.Update) (dsl string, table string, err error) {
	return "", "", errors.New("not support update")
}

func handleInsert(insert *sqlparser.Insert) (dsl string, table string, err error) {
	return "", "", errors.New("not support insert")
}

func handleDelete(delete *sqlparser.Delete) (dsl string, table string, err error) {
	return "", "", errors.New("not support delete")
}
