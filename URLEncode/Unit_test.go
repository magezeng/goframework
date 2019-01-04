package utils

import (
	"fmt"
	"testing"
)

type BClass struct {
	A string
}

func TestEncode(t *testing.T) {
	tempMap := map[string]interface{}{"a": 10, "b": 3.4, "c": 12.6}
	tempString := URLEncode(tempMap)
	fmt.Println(tempString)
}
