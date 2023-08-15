/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-15 09:55
**/

package test

import (
	"encoding/json"
	"github.com/lemonyxk/eutils/esql"
	"os"
	"reflect"
	"testing"
)

var testData map[string]interface{}

func init() {
	var f, err = os.Open("./data.json")
	if err != nil {
		panic(err)
	}
	defer func() { _ = f.Close() }()

	err = json.NewDecoder(f).Decode(&testData)
	if err != nil {
		panic(err)
	}
}

func TestAll(t *testing.T) {
	for k, v := range testData {
		var dsl, _, err = esql.Convert(k)
		if err != nil {
			t.Fatal(err)
		}

		var dslMap map[string]interface{}
		err = json.Unmarshal([]byte(dsl), &dslMap)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(dslMap, v) {
			t.Fatalf("sql: %s, want: %v, got: %v", k, v, dsl)
		}
	}
}
