package jsonpath

import (
	"testing"
	"encoding/json"
)

var jsonData = `
{
    "a": 3,
    "b": [
        1,
        2
    ],
    "c": {
        "d": 1
    },
    "e": {
        "f": {
            "c": 3
        }
    },
    "m": [
        {
            "n": 1
        }
    ]
}
`

func TestBaseSet(t *testing.T) {
	m := map[string]interface{}{}
	err := json.Unmarshal([]byte(jsonData), &m)
	if err != nil {
		t.Fatalf("数据json转换异常:%+v", err)
	}
	err = BaseSet(m, "f.h", ".", `["1","2"]`)
	if err != nil {
		t.Fatalf("数组数据设置异常:%+v", err)
	}
	err = BaseSet(m, "m.1.n", ".", "string")
	if err != nil {
		t.Fatalf("字符串数据设置异常:%+v", err)
	}
	err = BaseSet(m, "p", ".", `{"c": 3}`)
	if err != nil {
		t.Fatalf("对象数据设置异常:%+v", err)
	}
	mm := map[string]int{
		"a": 1,
	}
	err = BaseSet(m, "q.t", ".", mm)
	if err != nil {
		t.Fatalf("map数据设置异常:%+v", err)
	}
	arr := []int{1, 2}
	err = BaseSet(m, "u", ".", arr)
	if err != nil {
		t.Fatalf("Array数据设置异常:%+v", err)
	}
	b, _ := json.Marshal(m)
	t.Log(string(b))
}
