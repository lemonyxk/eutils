/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:43
**/

package esql

import "encoding/json"

type M map[string]interface{}

func (m M) String() string {
	bts, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(bts)
}

type A []interface{}
