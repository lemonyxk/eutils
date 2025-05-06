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

/*

SHOULD(`abc#2` like aaa) calc score [cuz no filter] : "should": [{ "bool": { "filter": [{ "term": { "city": { "value": "b" } } }] } }, { "bool": { "filter": [{ "term": { "city": { "value": "c" } } }] } }]
should not calc score : use filter eq sql `or` to handle : "should": [{ "bool": { "must": [{ "match_phrase": { "abc": { "boost": 2, "query": "aaa" } } }] } }]

*/

func main() {

	//var sql = "select a.b from a where SHOULD(`abc#2` BETWEEN 1 and 2, efg = 2) and `id#10` = 1 and `title#3` is not null and ((name = 'a' or name = 'b') and SHOULD(x = 1 , xx = 2) or c=2 and (age = 1 or age = 2)) and title like `1%1` order by id desc limit 10, 20"

	//var sql = "SELECT a,b FROM `social.post*` WHERE `categories_id` IN ('171') AND `crate_time` BETWEEN 1699027200 AND 1724947200 ORDER BY `CALC(income.purchase.aggregate.balance)` , SEARCH_AFTER(  '1' , 'a1'  ) DESC LIMIT 0,10"

	var sql = "SELECT * FROM `social.post*` WHERE SHOULD(`abc#2` like aaa) and name=a and (city = b or city = c)"

	log.Println(esql.Convert(sql))

}
