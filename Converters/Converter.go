package Converters

import (
	"encoding/json"
	"strconv"
)

/*
 int convert string
*/
func Int2String(arg int) string {
	return strconv.Itoa(arg)
}

/*
 int64 convert string
*/
func Int642String(arg int64) string {
	return strconv.FormatInt(int64(arg), 10)
}

/*
  json string convert struct
*/
func JsonString2Struct(jsonString string, result interface{}) error {
	jsonBytes := []byte(jsonString)
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
  json byte array convert struct
*/
func JsonBytes2Struct(jsonBytes []byte, result interface{}) error {
	err := json.Unmarshal(jsonBytes, result)
	return err
}

/*
 struct convert json string
*/
func Struct2JsonString(structt interface{}) (jsonString string, err error) {
	data, err := json.Marshal(structt)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

/*
  return decimalism string 9223372036854775807 -> "9223372036854775807"
*/
func Int64ToString(arg int64) string {
	return strconv.FormatInt(arg, 10)
}

func IntToString(arg int) string {
	return strconv.Itoa(arg)
}

func StringToInt64(arg string) int64 {
	value, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		return 0
	} else {
		return value
	}
}

func StringToInt(arg string) int {
	value, err := strconv.Atoi(arg)
	if err != nil {
		return 0
	} else {
		return value
	}
}