/**
* @program: es
*
* @description:
*
* @author: lemo
*
* @create: 2023-05-21 18:16
**/

package eutils

import (
	"reflect"
	"strings"
	"sync"
)

type M map[string]any

type Mapping struct {
	defaultKeyword bool
	textAsKeyword  bool
	longAsKeyword  bool
	ignoreAbove    int
	visited        map[uintptr]bool
	deep           int
	withTag        bool
	ignoreNil      bool
	mux            sync.Mutex
}

func (m *Mapping) DefaultKeyword(b bool) {
	m.defaultKeyword = b
}

func (m *Mapping) WithTag(b bool) {
	m.withTag = b
}

func (m *Mapping) IgnoreNil(b bool) {
	m.ignoreNil = b
}

func (m *Mapping) TextAsKeyword(b bool) {
	m.textAsKeyword = b
}

func (m *Mapping) LongAsKeyword(b bool) {
	m.longAsKeyword = b
}

func (m *Mapping) IgnoreAbove(i int) {
	if i < 0 {
		panic("ignore above must greater than 0")
	}

	if i > 32766 {
		panic("ignore above must less than 32766")
	}

	m.ignoreAbove = i
}

func NewMapping() *Mapping {
	return &Mapping{
		visited:        make(map[uintptr]bool),
		defaultKeyword: true,
		withTag:        false,
		ignoreNil:      true,
		textAsKeyword:  false,
		longAsKeyword:  false,
		ignoreAbove:    256,
		deep:           0,
	}
}

func (m *Mapping) GenerateMapping(t any) M {

	m.mux.Lock()
	defer m.mux.Unlock()

	var rv = reflect.ValueOf(t)

	var properties = M{}

	var mapping = M{
		"mappings": M{
			"properties": properties,
		},
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		panic("t must be struct")
	}

	m.printStruct(properties, "", rv)

	return mapping
}

func (m *Mapping) format(mapping map[string]any, key string, rv reflect.Value, tag reflect.StructTag) {
	switch rv.Kind() {

	// SIMPLE TYPE
	case reflect.Bool:
		// ignore
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Complex64, reflect.Complex128:
		// ignore
	case reflect.Float32, reflect.Float64:
		// ignore
	case reflect.String:
		// ignore
	case reflect.Func:
		// ignore
	case reflect.UnsafePointer:
		// ignore
	case reflect.Chan:
		// ignore
	case reflect.Invalid:
		// ignore

	// COMPLEX TYPE
	case reflect.Map:
		m.printMap(mapping, key, rv, tag)
	case reflect.Struct:
		m.printStruct(mapping, key, rv)
	case reflect.Array, reflect.Slice:
		m.printSlice(mapping, key, rv, tag)
	case reflect.Ptr:
		if rv.CanInterface() {
			m.printPtr(mapping, key, rv, tag)
		}
	case reflect.Interface:
		m.format(mapping, key, rv.Elem(), tag)
	default:
		// ignore
	}
}

func (m *Mapping) printMap(mapping map[string]any, key string, v reflect.Value, tag reflect.StructTag) {

	var d = m.deep
	m.deep++

	if v.Len() == 0 {
		// ignore
		// cuz you don't know the type of key and value
		m.deep = d
		return
	}

	if m.visited[v.Pointer()] {
		// ignore
		// repeat reference
		m.deep = d
		return
	}

	m.visited[v.Pointer()] = true

	var newMapping = M{}
	mapping[key] = M{
		"properties": newMapping,
	}

	keys := v.MapKeys()
	for i := 0; i < v.Len(); i++ {
		value := v.MapIndex(keys[i])
		var fieldName = keys[i].String()

		var name, parse = m.parseElasticTag(tag)

		if parse == nil && m.withTag {
			continue
		}

		if parse == nil {
			parse = &parser{}
		}

		if parse.Ignore {
			continue
		}

		if name == "" {
			name = fieldName
		}

		var defaultType = m.defaultType(value)

		var tp = parse.Type

		if parse.Type == "" {
			tp = defaultType
		}

		var t = M{
			"type": tp,
		}

		if parse.Type == "keyword" {
			t["ignore_above"] = m.ignoreAbove
		}

		if parse.Analyzer != "" && tp == "text" {
			t["analyzer"] = parse.Analyzer
		}

		if parse.Keyword && tp == "text" && !m.textAsKeyword {
			t["fields"] = M{
				"keyword": M{
					"type":         "keyword",
					"ignore_above": m.ignoreAbove,
				},
			}
		}

		if m.textAsKeyword && tp == "text" && (parse.Type == "" && parse.Analyzer == "") {
			t["type"] = "keyword"
			t["ignore_above"] = m.ignoreAbove
		}

		if m.longAsKeyword && tp == "long" && (parse.Type == "" && parse.Analyzer == "") {
			t["type"] = "keyword"
			t["ignore_above"] = m.ignoreAbove
		}

		if parse.Index != nil {
			t["index"] = *parse.Index
		}

		newMapping[fieldName] = t

		m.format(newMapping, fieldName, value, tag)
	}

	m.deep = d
}

func (m *Mapping) printStruct(mapping map[string]any, key string, v reflect.Value) {

	var d = m.deep
	m.deep++

	if v.NumField() == 0 {
		// ignore
		// cuz you don't know the type of key and value
		m.deep = d
		return
	}

	var newMapping = M{}
	mapping[key] = M{
		"properties": newMapping,
	}

	for i := 0; i < v.NumField(); i++ {
		var field = v.Type().Field(i)

		fieldName := field.Name
		value := v.Field(i)

		if value.CanInterface() {

			// elasticsearch ignore anonymous field
			if field.Anonymous {
				m.format(mapping, key, value, field.Tag)
				continue
			}

			var name, parse = m.parseElasticTag(field.Tag)

			if parse == nil && m.withTag {
				continue
			}

			if parse == nil {
				parse = &parser{}
			}

			if parse.Ignore {
				continue
			}

			if name == "" {
				name = fieldName
			}

			var defaultType = m.defaultType(value)

			var tp = parse.Type

			if parse.Type == "" {
				tp = defaultType
			}

			var t = M{
				"type": tp,
			}

			if parse.Type == "keyword" {
				t["ignore_above"] = m.ignoreAbove
			}

			if parse.Analyzer != "" && tp == "text" {
				t["analyzer"] = parse.Analyzer
			}

			if parse.Keyword && tp == "text" && !m.textAsKeyword {
				t["fields"] = M{
					"keyword": M{
						"type":         "keyword",
						"ignore_above": m.ignoreAbove,
					},
				}
			}

			if m.textAsKeyword && tp == "text" && (parse.Type == "" && parse.Analyzer == "") {
				t["type"] = "keyword"
				t["ignore_above"] = m.ignoreAbove
			}

			if m.longAsKeyword && tp == "long" && (parse.Type == "" && parse.Analyzer == "") {
				t["type"] = "keyword"
				t["ignore_above"] = m.ignoreAbove
			}

			if parse.Index != nil {
				t["index"] = *parse.Index
			}

			// which type of value can be nil
			switch value.Type().Kind() {
			case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface:
				if value.Type().Kind() == reflect.Ptr && !value.CanInterface() {
					continue
				} else {
					if value.IsNil() {
						if m.ignoreNil {
							continue
						} else {
							if value.Type().Kind() == reflect.Interface {
								// you can not know the type of interface if it is nil
								continue
							} else {
								value = reflect.New(value.Type().Elem())
							}
						}
					}
				}
			}

			// first struct
			if key == "" {
				delete(mapping, key)
				mapping[name] = t
				m.format(mapping, name, value, field.Tag)
				continue
			}

			newMapping[name] = t

			// printTags(newMapping, name, value)
			m.format(newMapping, name, value, field.Tag)
		}
	}

	m.deep = d
}

func (m *Mapping) printSlice(mapping map[string]any, key string, v reflect.Value, tag reflect.StructTag) {

	var d = m.deep
	m.deep++

	if v.Len() == 0 {
		// ignore
		// cuz you don't know the type of key and value
		m.deep = d
		return
	}

	//  if is array, will be handled in printPtr
	if v.Kind() == reflect.Slice {
		if m.visited[v.Pointer()] {
			// repeat reference
			m.deep = d
			return
		}
		m.visited[v.Pointer()] = true
	}

	// only print first element
	// cuz elastic don't support array of different type
	// always use first element type no matter what type it is
	// example: []interface{}{1, "2", true} [[[1,2],3]]
	// but you have to make sure all elements are the same type
	// otherwise you will get error from elastic
	if v.Len() > 0 {
		m.format(mapping, key, v.Index(0), tag)
	} else {
		m.format(mapping, key, reflect.New(v.Type().Elem()), tag)
	}

	m.deep = d
}

func (m *Mapping) printPtr(mapping map[string]any, key string, v reflect.Value, tag reflect.StructTag) {

	if m.visited[v.Pointer()] {
		// repeat reference
		return
	}

	if v.Pointer() != 0 {
		m.visited[v.Pointer()] = true
	}

	if v.Elem().IsValid() {
		m.format(mapping, key, v.Elem(), tag)
	}
}

type parser struct {
	Keyword  bool
	Analyzer string
	Type     string
	Index    *bool
	Ignore   bool
}

func (m *Mapping) parseElasticTag(tag reflect.StructTag) (string, *parser) {
	var str = tag.Get("es")
	var js = tag.Get("json")
	var name = strings.Split(js, ",")[0]
	var tp string
	var keyword bool
	var analyzer string
	var ignore bool
	var index *bool
	var arr = strings.Split(str, ",")
	if (str == "" || len(arr) == 0) && m.withTag {
		return name, nil
	}

	if str == "-" {
		return name, nil
	}

	if js == "-" && m.withTag {
		return name, nil
	}

	for i := 0; i < len(arr); i++ {
		if strings.HasPrefix(arr[i], "type:") {
			tp = strings.TrimPrefix(arr[i], "type:")
		}
		if strings.HasPrefix(arr[i], "keyword:") {
			keyword = strings.TrimPrefix(arr[i], "keyword:") == "true"
		}
		if strings.HasPrefix(arr[i], "analyzer:") {
			analyzer = strings.TrimPrefix(arr[i], "analyzer:")
		}
		if strings.HasPrefix(arr[i], "index:") {
			var b = strings.TrimPrefix(arr[i], "index:") == "true"
			index = &b
		}
		if strings.HasPrefix(arr[i], "ignore:") {
			ignore = strings.TrimPrefix(arr[i], "ignore:") == "true"
		}
	}

	if !strings.Contains(str, "keyword") && m.defaultKeyword {
		keyword = true
	}

	return name, &parser{
		Keyword:  keyword,
		Analyzer: analyzer,
		Type:     tp,
		Index:    index,
		Ignore:   ignore,
	}
}

func (m *Mapping) defaultType(tp reflect.Value) string {
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	switch tp.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Complex64, reflect.Complex128:
		return "long"
	case reflect.Float32, reflect.Float64:
		return "double"
	case reflect.String:
		return "text"
	case reflect.Slice, reflect.Array:
		if tp.Type().Elem().Kind() == reflect.Uint8 {
			return "binary"
		}
		if tp.Len() > 0 {
			return m.defaultType(tp.Index(0))
		} else {
			return m.defaultType(reflect.New(tp.Type().Elem()))
		}
	default:
		return "text"
	}
}
