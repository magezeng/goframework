package Converters

import (
	"encoding/json"
	"reflect"
	"strconv"
)

/*
 int convert string
*/
func IntToString(arg int) string {
	return strconv.Itoa(arg)
}

/*
 int64 convert string
*/
func Int64ToString(arg int64) string {
	return strconv.FormatInt(int64(arg), 10)
}

/*
  json string convert struct
*/
func JsonStringToStruct(jsonString string, result interface{}) error {
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
  json byte array convert struct
*/
func JsonBytesToStruct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
 struct convert json string
*/
func StructToJsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func StringToInt64(arg string) int64 {
	value, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StructToMap(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func StringToMap(obj string) map[string]interface{} {
	var data = make(map[string]interface{})
	err := json.Unmarshal([]byte(obj), &data)
	if err != nil {
	}
	return data
}

func Float64ToString(arg float64) string {
	return strconv.FormatFloat(arg, 'f', -1, 64)
}

func StringToInt(arg string) int {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func SimpleInterfaceToString(arg interface{}) string {
	switch arg.(type) {
	case int:
		return IntToString(arg.(int))
	case int64:
		return Int64ToString(arg.(int64))
	case float64:
		return Float64ToString(arg.(float64))
	case string:
		return arg.(string)
	default:
		return ""
	}
}
