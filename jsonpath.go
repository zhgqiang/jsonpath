package jsonpath

import (
	"strings"
	"encoding/json"
	"regexp"
	"strconv"
)

var aa = `
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

const (
	STRING  = "string"
	INT     = "int"
	FLOAT64 = "float64"
	FLOAT32 = "float32"
	MAP     = "map"
	ARRAY   = "array"
	DEFAULT = "default"
)

// m 数据存储的json格式生成
// path 转换路径 "m.0.n"
// pathSplit path的分割符
// value 数据的值,map、array、string、number等
// 例如：var m map[string]interface{} = map[string]interface{}{}
// BaseSet(m, "m.0.n", ".", 3)
func BaseSet(m map[string]interface{}, path, pathSplit string, value interface{}) error {
	if m == nil {
		m = map[string]interface{}{}
	}
	if strings.TrimSpace(pathSplit) == "" {
		pathSplit = "."
	}
	if strings.Contains(path, pathSplit) {
		pathArray := strings.Split(path, pathSplit)
		return setValue(0, pathArray, value, m, m)
	} else {
		return setMapValue(m, path, value)
	}
}

func setValue(index int, arr []string, value interface{}, preData interface{}, nextData interface{}) error {
	match, err := regexp.MatchString("^[0-9]\\d*$", arr[index])
	if err != nil {
		return err
	}
	if match {
		i, err := strconv.Atoi(arr[index])
		if err != nil {
			return err
		}
		m := nextData.([]interface{})

		if index == len(arr)-1 {
			return setArrayValue(m, i, value)
		}
		if len(m) <= i {
			nn := make([]interface{}, i+1)
			for ii := 0; ii < len(m); ii++ {
				nn[ii] = m[ii]
			}
			m = nn
		}
		val := m[i]

		if val == nil {
			match, err := regexp.MatchString("^[0-9]\\d*$", arr[index+1])
			if err != nil {
				return err
			}

			if match {
				j, err := strconv.Atoi(arr[index+1])
				if err != nil {
					return err
				}
				val = make([]interface{}, j+1)
			} else {
				val = map[string]interface{}{}
			}

			m[i] = val
		}
		match1, err := regexp.MatchString("^[0-9]\\d*$", arr[index-1])
		if err != nil {
			return err
		}

		if match1 {
			mPre := preData.([]interface{})
			j, err := strconv.Atoi(arr[index-1])
			if err != nil {
				return err
			}
			mPre[j] = m
		} else {
			mPre := preData.(map[string]interface{})
			mPre[arr[index-1]] = m
		}
		index++
		return setValue(index, arr, value, m, val)
	} else {
		m := nextData.(map[string]interface{})
		if index == len(arr)-1 {
			return setMapValue(m, arr[index], value)
		}
		val, ok := m[arr[index]]
		if !ok {
			match, err := regexp.MatchString("^[0-9]\\d*$", arr[index+1])
			if err != nil {
				return err
			}

			if match {
				j, err := strconv.Atoi(arr[index+1])
				if err != nil {
					return err
				}
				val = make([]interface{}, j+1)
			} else {
				val = map[string]interface{}{}
			}
			m[arr[index]] = val
		}
		index++
		return setValue(index, arr, value, m, val)
	}
}

func setMapValue(m map[string]interface{}, key string, value interface{}) error {
	if jType(value) == STRING {
		if valiateDataType(value.(string)) == MAP {
			mm := map[string]interface{}{}
			err := json.Unmarshal([]byte(value.(string)), &mm)
			if err != nil {
				return err
			}
			m[key] = mm
			return nil
		} else if valiateDataType(value.(string)) == ARRAY {
			mm := []interface{}{}
			err := json.Unmarshal([]byte(value.(string)), &mm)
			if err != nil {
				return err
			}
			m[key] = mm
			return nil
		} else {
			m[key] = value
			return nil
		}
	} else if jType(value) == MAP {
		m[key] = value.(map[string]interface{})
		return nil
	} else if jType(value) == ARRAY {
		m[key] = value.([]interface{})
		return nil
	} else {
		m[key] = value
		return nil
	}
}

func setArrayValue(m []interface{}, i int, value interface{}) error {
	if jType(value) == STRING {
		if valiateDataType(value.(string)) == MAP {
			mm := map[string]interface{}{}
			err := json.Unmarshal([]byte(value.(string)), &mm)
			if err != nil {
				return err
			}
			m[i] = mm
			return nil
		} else if valiateDataType(value.(string)) == ARRAY {
			mm := []interface{}{}
			err := json.Unmarshal([]byte(value.(string)), &mm)
			if err != nil {
				return err
			}
			m[i] = mm
			return nil
		} else {
			m[i] = value
			return nil
		}

	} else if jType(value) == MAP {
		m[i] = value.(map[string]interface{})
		return nil
	} else if jType(value) == ARRAY {
		m[i] = value.([]interface{})
		return nil
	} else {
		m[i] = value
		return nil
	}
}

func valiateDataType(value string) string {
	mmm := []interface{}{}
	err := json.Unmarshal([]byte(value), &mmm)

	if err != nil {
		mm := map[string]interface{}{}
		err = json.Unmarshal([]byte(value), &mm)
		if err != nil {
			return STRING
		}
		if len(mm) != 0 {
			return MAP
		}
	}

	if len(mmm) != 0 {
		return ARRAY
	}
	mm := map[string]interface{}{}
	err = json.Unmarshal([]byte(value), &mm)
	if err != nil {
		return STRING
	}
	if len(mm) != 0 {
		return MAP
	}

	return STRING
}

func jType(i interface{}) string { //函数t 有一个参数i
	switch i.(type) { //多选语句switch
	case string:
		return STRING
	case int:
		return INT
	case float64:
		return FLOAT64
	case float32:
		return FLOAT32
	case map[string]interface{}:
		return MAP
	case []interface{}:
		return ARRAY
	default:
		return DEFAULT
	}
}
