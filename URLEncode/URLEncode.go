package utils

import (
	"bytes"
	. "net/url"
	"sort"
	"strconv"
)

func URLEncode(tempMap map[string]interface{}) string {
	values := MapToValues(tempMap)
	return Encode(values)
}

func toString(value interface{}) (stringValue string) {
	switch value.(type) {
	default:
		stringValue = "不支持此类型"
	case bool:
		if value.(bool) {
			stringValue = "true"
		} else {
			stringValue = "false"
		}
	case int:
		stringValue = string(value.(int))
	case float32:
		stringValue = strconv.FormatFloat(float64(value.(float32)), 'f', 20, 64)
	case float64:
		stringValue = strconv.FormatFloat(value.(float64), 'f', 20, 64)
	}
	return
}

func MapToValues(tempMap map[string]interface{}) Values {
	tempValues := Values{}
	for k, v := range tempMap {
		tempValues.Add(k, toString(v))
	}
	return tempValues
}

func Encode(v Values) string {
	if v == nil {
		return ""
	}
	var buf bytes.Buffer
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		prefix := QueryEscape(k) + "="
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(prefix)
			buf.WriteString(QueryEscape(v))
		}
	}
	return buf.String()
}
