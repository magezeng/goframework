package ArrValueTranslate

import (
	"fmt"
	"testing"
)

type BClass struct {
	A string
}

func TestArrValueTranslate(t *testing.T) {

	var fromArray []interface{}
	fromArray = []interface{}{BClass{"123"}, BClass{"231"}}

	var toArray = []BClass{}

	ArrValueTranslate(fromArray, &toArray)
	fmt.Println(toArray)
}
