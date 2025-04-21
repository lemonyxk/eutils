/**
* @program: engine
*
* @create: 2025-04-18 13:50
**/

package script

import (
	"bytes"
	"context"
	"github.com/elastic/go-elasticsearch/v9"
	"github.com/elastic/go-elasticsearch/v9/esapi"
	"github.com/lemonyxk/kitty/json"
	"github.com/lemonyxk/kitty/kitty"
)

/*

POST test_1/_update/1
{
  "script": {
    "id":"update",
    "params": {
      "$remove":{
        "target":{"opp":{"1.bb":1}},
        "index":0
      },
      "$add":{
        "target":{"opp":{"1.bb":1}},
        "value":1 or [1,2,3]
      },
      "$pull":{
        "target":{"opp":{"1.bb":1}},
        "value":1 or [1,2,3]
      },
      "$replace":{
        "target":{"opp":{"1.bb":1}},
        "old":1 or [1,2,3]
		"new":2
      },
	  "$unset":{
		"target":{"opp":{"1.bb":1}}
	  },
	  "$inc":{
		"opp":{"1.bb":1}
	  },
	  "$set":{
		"opp":{"1.bb":1}
	  }
    }
  }
}

*/

func CreateUpdateScript(client *elasticsearch.Client) {

	var script = kitty.M{
		"script": kitty.M{
			"source": UpdateScript,
			"lang":   "painless",
		},
	}

	var bts, err = json.Marshal(script)
	if err != nil {
		panic(err)
	}

	var req = esapi.PutScriptRequest{
		ScriptID: "update",
		Body:     bytes.NewReader(bts),
	}

	res, err := req.Do(context.Background(), client)
	if err != nil {
		panic(err)
	}

	defer func() { _ = res.Body.Close() }()

	if res.IsError() {
		panic(res.String())
	}

	defer func() { _ = res.Body.Close() }()

	println(res.String())
}

var UpdateScript = `
// 递归函数，构建路径
String buildPath(Map m) {
    if (m.size() == 0) {
        return "";
    }
    def it = m.entrySet().iterator();
    def entry = it.next();
    String key = entry.getKey();
    def val = entry.getValue();
    if (val instanceof Map) {
        return key + "." + buildPath(val);
    } else {
        return key;
    }
}

// 判断字符串是否是数字（整数）
boolean isInteger(String s) {
	if (!(s.startsWith("[") && s.endsWith("]"))) {
		return false;
	}
    try {
       Integer.parseInt(s.substring(1, s.length() - 1));
       return true;
    } catch (Exception e) {
       return false;
    }
}

int parseInt(String s) {
	return Integer.parseInt(s.substring(1, s.length() - 1));
}

Map flattenMap(Map input, String prefix) {
    Map result = new HashMap();
    for (entry in input.entrySet()) {
        String key = entry.getKey();
        def val = entry.getValue();
        String newKey = prefix == null || prefix == "" ? key : prefix + "." + key;
        if (val instanceof Map) {
            Map subMap = flattenMap(val, newKey);
            for (subEntry in subMap.entrySet()) {
                result.put(subEntry.getKey(), subEntry.getValue());
            }
        } else {
            result.put(newKey, val);
        }
    }
    return result;
}

// 递归自增函数，支持数组索引路径，路径不存在则创建对应结构
void incrementField(def source, String[] path, int index, def incValue) {
    String key = path[index];
    boolean last = (index == path.length - 1);

    if (isInteger(key)) {
        int idx = parseInt(key);
        // 如果当前层不是List，初始化为List
        if (!(source instanceof List)) {
            // 注意：这里不能直接赋值给source，因为Java参数是值传递，需要调用方确保传入的是可变引用
            // 但Painless中ctx._source是引用类型，通常没问题
            // 这里假设source是引用，直接转换
            // 如果你调用时发现问题，可以调整调用逻辑
            source = new ArrayList();
        }
        List list = (List) source;
        // 扩展数组长度，保证索引有效
        while (list.size() <= idx) {
            list.add(null);
        }

        if (last) {
            // 最后一层，执行自增
            if (list.get(idx) instanceof Number) {
                list.set(idx, list.get(idx) + incValue);
            } else {
                list.set(idx, incValue);
            }
        } else {
            def next = list.get(idx);
            String nextKey = path[index + 1];
            boolean nextIsIndex = isInteger(nextKey);

            if (next == null) {
                // 根据下一级路径判断创建 Map 或 List
                next = nextIsIndex ? new ArrayList() : new HashMap();
                list.set(idx, next);
            } else if (!(next instanceof Map) && !(next instanceof List)) {
                // 如果类型不匹配，覆盖为合适类型
                next = nextIsIndex ? new ArrayList() : new HashMap();
                list.set(idx, next);
            }
            incrementField(next, path, index + 1, incValue);
        }
    } else {
        // key 是 Map 的 key
        if (!(source instanceof Map)) {
            source = new HashMap();
        }
        Map map = (Map) source;

        if (last) {
            if (map.containsKey(key) && map.get(key) instanceof Number) {
                map.put(key, map.get(key) + incValue);
            } else {
                map.put(key, incValue);
            }
        } else {
            def next = map.get(key);
            String nextKey = path[index + 1];
            boolean nextIsIndex = isInteger(nextKey);

            if (next == null) {
                next = nextIsIndex ? new ArrayList() : new HashMap();
                map.put(key, next);
            } else if (!(next instanceof Map) && !(next instanceof List)) {
                next = nextIsIndex ? new ArrayList() : new HashMap();
                map.put(key, next);
            }
            incrementField(next, path, index + 1, incValue);
        }
    }
}

void processIncrements(Map source, Map increments) {
    Map flatMap = flattenMap(increments, "");
    for (entry in flatMap.entrySet()) {
        String targetPath = entry.getKey();
        def value = entry.getValue();

		String[] path = targetPath.splitOnToken(".");

        incrementField(source, path, 0, value);
    }
}


// 递归自增函数，支持数组索引路径，路径不存在则创建对应结构
void setField(def source, String[] path, int index, def incValue) {
    String key = path[index];
    boolean last = (index == path.length - 1);

    if (isInteger(key)) {
        int idx = parseInt(key);
        // 如果当前层不是List，初始化为List
        if (!(source instanceof List)) {
            // 注意：这里不能直接赋值给source，因为Java参数是值传递，需要调用方确保传入的是可变引用
            // 但Painless中ctx._source是引用类型，通常没问题
            // 这里假设source是引用，直接转换
            // 如果你调用时发现问题，可以调整调用逻辑
            source = new ArrayList();
        }
        List list = (List) source;
        // 扩展数组长度，保证索引有效
        while (list.size() <= idx) {
            list.add(null);
        }

        if (last) {
            list.set(idx, incValue);
        } else {
            def next = list.get(idx);
            String nextKey = path[index + 1];
            boolean nextIsIndex = isInteger(nextKey);

            if (next == null) {
                // 根据下一级路径判断创建 Map 或 List
                next = nextIsIndex ? new ArrayList() : new HashMap();
                list.set(idx, next);
            } else if (!(next instanceof Map) && !(next instanceof List)) {
                // 如果类型不匹配，覆盖为合适类型
                next = nextIsIndex ? new ArrayList() : new HashMap();
                list.set(idx, next);
            }
            setField(next, path, index + 1, incValue);
        }
    } else {
        // key 是 Map 的 key
        if (!(source instanceof Map)) {
            source = new HashMap();
        }
        Map map = (Map) source;

        if (last) {
            map.put(key, incValue);
        } else {
            def next = map.get(key);
            String nextKey = path[index + 1];
            boolean nextIsIndex = isInteger(nextKey);

            if (next == null) {
                next = nextIsIndex ? new ArrayList() : new HashMap();
                map.put(key, next);
            } else if (!(next instanceof Map) && !(next instanceof List)) {
                next = nextIsIndex ? new ArrayList() : new HashMap();
                map.put(key, next);
            }
            setField(next, path, index + 1, incValue);
        }
    }
}

void processSet(Map source, Map sets) {
    Map flatMap = flattenMap(sets, "");
    for (entry in flatMap.entrySet()) {
        String targetPath = entry.getKey();
        def value = entry.getValue();

		String[] path = targetPath.splitOnToken(".");

        setField(source, path, 0, value);
    }
}

// 只查询路径，路径不存在返回null，支持数组索引访问
def getParentAndKeyIfExists(def source, String[] path) {
    def current = source;
    for (int i = 0; i < path.length - 1; i++) {
        String key = path[i];
        if (isInteger(key)) {
            int idx = parseInt(key);
            if (!(current instanceof List)) {
                return null;
            }
            List list = (List) current;
            if (idx >= list.size()) {
                return null;
            }
            current = list.get(idx);
            if (current == null) {
                return null;
            }
        } else {
            if (!(current instanceof Map)) {
                return null;
            }
            Map map = (Map) current;
            if (!map.containsKey(key)) {
                return null;
            }
            current = map.get(key);
            if (current == null) {
                return null;
            }
        }
    }
    // 返回父容器和最后一级key
    return ['parent': current, 'key': path[path.length - 1]];
}

def getOrCreate(def current, String[] path, int index) {
    if (index == path.length - 1) {
        return ['parent': current, 'key': path[index]];
    }
    String key = path[index];
    String nextKey = path[index + 1];
    boolean nextIsIndex = isInteger(nextKey);

    if (isInteger(key)) {
        int idx = parseInt(key);
        if (!(current instanceof List)) {
            throw new IllegalArgumentException("Expected List at " + key);
        }
        List list = (List) current;
        while (list.size() <= idx) {
            list.add(null);
        }
        if (list.get(idx) == null) {
            list.set(idx, nextIsIndex ? new ArrayList() : new HashMap());
        }
        return getOrCreate(list.get(idx), path, index + 1);
    } else {
        if (!(current instanceof Map)) {
            throw new IllegalArgumentException("Expected Map at " + key);
        }
        Map map = (Map) current;
        if (!map.containsKey(key) || map.get(key) == null) {
            map.put(key, nextIsIndex ? new ArrayList() : new HashMap());
        }
        return getOrCreate(map.get(key), path, index + 1);
    }
}

// $unset：删除指定路径字段，路径不存在则忽略
void unsetFieldIfExists(def source, String targetPath) {
	String[] path = targetPath.splitOnToken(".");

    def containerAndKey = getParentAndKeyIfExists(source, path);
    if (containerAndKey == null) {
        return; // 路径不存在，忽略
    }
    def parent = containerAndKey['parent'];
    String key = containerAndKey['key'];

    if (isInteger(key)) {
        int idx = parseInt(key);
        if (!(parent instanceof List)) {
            return;
        }
        List list = (List) parent;
        if (idx >= list.size()) {
            return;
        }
        // 删除数组中指定索引元素，使用remove
        list.remove(idx);
    } else {
        if (!(parent instanceof Map)) {
            return;
        }
        Map map = (Map) parent;
        if (map.containsKey(key)) {
            map.remove(key);
        }
    }
}

// 替换数组中所有等于old数组中任一值的元素为new
void replaceInArrayIfExists(Map source, String targetPath, List oldVals, def newVal) {
	String[] path = targetPath.splitOnToken(".");

    Map containerAndKey = getParentAndKeyIfExists(source, path);
    if (containerAndKey == null) return;
    Map parent = containerAndKey['parent'];
    String key = containerAndKey['key'];
    if (parent.containsKey(key) && parent[key] instanceof List) {
        List arr = parent[key];
        for (int i = 0; i < arr.size(); i++) {
            for (int j = 0; j < oldVals.size(); j++) {
                if (arr[i] == oldVals[j]) {
                    arr[i] = newVal;
                    break;
                }
            }
        }
    }
}

// $pull：删除数组中所有等于value数组中任一元素的元素
void pullFromArrayIfExists(Map source, String targetPath, List values) {
	String[] path = targetPath.splitOnToken(".");

    Map containerAndKey = getParentAndKeyIfExists(source, path);
    if (containerAndKey == null) return;
    Map parent = containerAndKey['parent'];
    String key = containerAndKey['key'];
    if (parent.containsKey(key) && parent[key] instanceof List) {
        List arr = parent[key];
        for (int i = arr.size() - 1; i >= 0; i--) {
            for (int j = 0; j < values.size(); j++) {
                if (arr[i] == values[j]) {
                    arr.remove(i);
                    break;
                }
            }
        }
    }
}

void addToArrayOrCreate(def source, String targetPath, List values) {
	String[] path = targetPath.splitOnToken(".");

    def containerAndKey = getOrCreate(source, path, 0);
    def parent = containerAndKey['parent'];
    String key = containerAndKey['key'];

    if (isInteger(key)) {
        int idx = parseInt(key);
        if (!(parent instanceof List)) {
            throw new IllegalArgumentException("Expected List at final segment");
        }
        List list = (List) parent;
        while (list.size() <= idx) {
            list.add(null);
        }
        if (!(list.get(idx) instanceof List)) {
            list.set(idx, new ArrayList());
        }
        List arr = (List) list.get(idx);
        for (int i = 0; i < values.size(); i++) {
            arr.add(values.get(i));
        }
    } else {
        if (!(parent instanceof Map)) {
            throw new IllegalArgumentException("Expected Map at final segment");
        }
        Map map = (Map) parent;
        if (!map.containsKey(key) || !(map.get(key) instanceof List)) {
            map.put(key, new ArrayList());
        }
        List arr = (List) map.get(key);
        for (int i = 0; i < values.size(); i++) {
            arr.add(values.get(i));
        }
    }
}

// $remove 函数：删除指定路径数组中指定索引元素，路径或索引不存在则忽略
void removeIndexFromArrayIfExists(def source, String targetPath, int indexToRemove) {
	String[] path = targetPath.splitOnToken(".");

    def containerAndKey = getParentAndKeyIfExists(source, path);
    if (containerAndKey == null) {
        return; // 路径不存在，忽略
    }
    def parent = containerAndKey['parent'];
    String key = (String) containerAndKey['key'];

    if (isInteger(key)) {
        int idx = parseInt(key);
        if (!(parent instanceof List)) {
            return; // 父不是数组，忽略
        }
        List list = (List) parent;
        if (idx >= list.size()) {
            return; // 索引越界，忽略
        }
        def arr = list.get(idx);
        if (!(arr instanceof List)) {
            return; // 目标不是数组，忽略
        }
        List arrList = (List) arr;
        if (indexToRemove >= 0 && indexToRemove < arrList.size()) {
            arrList.remove(indexToRemove);
        }
    } else {
        if (!(parent instanceof Map)) {
            return; // 父不是对象，忽略
        }
        Map map = (Map) parent;
        if (!map.containsKey(key)) {
            return; // 字段不存在，忽略
        }
        if (!(map.get(key) instanceof List)) {
            return; // 目标不是数组，忽略
        }
        List arrList = (List) map.get(key);
        if (indexToRemove >= 0 && indexToRemove < arrList.size()) {
            arrList.remove(indexToRemove);
        }
    }
}

if (params.containsKey('$set')) {
    def setParam = params['$set'];
    processSet(ctx._source, setParam);
}

if (params.containsKey('$inc')) {
    def incParam = params['$inc'];
    processIncrements(ctx._source, incParam);
}

if (params.containsKey('$unset')) {
    def unsetParam = params['$unset'];
    if (unsetParam.containsKey('target')) {
        String targetPath;
        if (unsetParam.target instanceof Map) {
            targetPath = buildPath(unsetParam.target);
        } else {
            targetPath = unsetParam.target;
        }
        unsetFieldIfExists(ctx._source, targetPath);
    }
}

if (params.containsKey('$replace')) {
    def replaceParam = params['$replace'];
    if (replaceParam.containsKey('target') && replaceParam.containsKey('old') && replaceParam.containsKey('new')) {
        String targetPath;
        if (replaceParam.target instanceof Map) {
            targetPath = buildPath(replaceParam.target);
        } else {
            targetPath = replaceParam.target;
        }
        List oldList = new ArrayList();
        if (replaceParam.old instanceof List) {
            oldList = replaceParam.old;
        } else {
            oldList.add(replaceParam.old);
        }
        replaceInArrayIfExists(ctx._source, targetPath, oldList, replaceParam.new);
    }
}

if (params.containsKey('$pull')) {
    def pullParam = params['$pull'];
    if (pullParam.containsKey('target') && pullParam.containsKey('value')) {
        String targetPath;
        if (pullParam.target instanceof Map) {
            targetPath = buildPath(pullParam.target);
        } else {
            targetPath = pullParam.target;
        }
        List values = new ArrayList();
        if (pullParam.value instanceof List) {
            values = pullParam.value;
        } else {
            values.add(pullParam.value);
        }
        pullFromArrayIfExists(ctx._source, targetPath, values);
    }
}

if (params.containsKey('$add')) {
    def addParam = params['$add'];
    if (addParam.containsKey('target') && addParam.containsKey('value')) {
        String targetPath;
        if (addParam.target instanceof Map) {
            targetPath = buildPath(addParam.target);
        } else {
            targetPath = addParam.target;
        }
        List values = new ArrayList();
        if (addParam.value instanceof List) {
            values = addParam.value;
        } else {
            values.add(addParam.value);
        }
        addToArrayOrCreate(ctx._source, targetPath, values);
    }
}

if (params.containsKey('$remove')) {
    def removeParam = params['$remove'];
    if (removeParam.containsKey('target') && removeParam.containsKey('index')) {
        String targetPath;
        if (removeParam.target instanceof Map) {
            targetPath = buildPath(removeParam.target);
        } else {
            targetPath = removeParam.target;
        }
        removeIndexFromArrayIfExists(ctx._source, targetPath, removeParam.index);
    }
}
`
