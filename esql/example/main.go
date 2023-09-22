/**
* @program: elasticsql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 18:50
**/

package main

import (
	"fmt"
	"github.com/lemonyxk/eutils/esql"
)

func main() {

	var sql = "select includes(TITLE),   excludes(`DESC`) from a where abc = 1 and id = 1 and ((name = 'a' or name = 'b') or c=2 and (age = 1 or age = 2)) order by id desc limit 10, 20"

	fmt.Println(esql.Convert(sql))

}
