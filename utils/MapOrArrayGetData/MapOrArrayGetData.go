package MapOrArrayGetData

import (
	"errors"
	"fmt"
)

type mapOrArrayOrSub struct {
	data        interface{}
	accessDepth int
}

func New(data interface{}) mapOrArrayOrSub {
	return mapOrArrayOrSub{data, 0}
}

func (parent mapOrArrayOrSub) GetData() (data interface{}, err error) {
	if tempErr, isError := parent.data.(error); isError {
		err = tempErr
		return
	}

	data = parent.data
	return
}

func (parent mapOrArrayOrSub) GetSubWithKey(key string) mapOrArrayOrSub {

	if _, isError := parent.data.(error); isError {
		return parent
	}

	if tempMap, ok := parent.data.(map[string]interface{}); ok {
		if tempValue, haveKey := tempMap[key]; haveKey {
			return mapOrArrayOrSub{tempValue, parent.accessDepth + 1}
		} else {
			fmt.Println("MapOrArrayOrSub:访问深度" + fmt.Sprint(parent.accessDepth))
			return mapOrArrayOrSub{errors.New("该Map没有Key" + fmt.Sprint(key)), 0}
		}
	} else {
		fmt.Println("MapOrArrayOrSub:访问深度" + fmt.Sprint(parent.accessDepth))
		return mapOrArrayOrSub{errors.New("该对象不是一个Map"), 0}
	}
}

func (parent mapOrArrayOrSub) GetSubWithIndex(index int) mapOrArrayOrSub {

	if _, isError := parent.data.(error); isError {
		return parent
	}

	if tempArr, ok := parent.data.([]interface{}); ok {
		if len(tempArr) < index+1 {
			fmt.Println("MapOrArrayOrSub:访问深度" + fmt.Sprint(parent.accessDepth))
			return mapOrArrayOrSub{errors.New("超出该数组长度"), 0}
		}
		return mapOrArrayOrSub{tempArr[index], parent.accessDepth + 1}
	} else {
		fmt.Println("MapOrArrayOrSub:访问深度" + fmt.Sprint(parent.accessDepth))
		return mapOrArrayOrSub{errors.New("该对象不是一个Array"), 0}
	}
}
