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
	"github.com/lemonyxk/eutils/esql"
	"log"
)

func main() {

	var sql = "select a.b from a where SHOULD(`abc#2` = 1, efg = 2) and `id#10` = 1 and `title#3` is not null and ((name = 'a' or name = 'b') and SHOULD(x = 1 , xx = 2) or c=2 and (age = 1 or age = 2)) order by id desc limit 10, 20"

	log.Println(esql.Convert(sql))

}
