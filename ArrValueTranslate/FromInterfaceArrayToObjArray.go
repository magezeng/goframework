package ArrValueTranslate

import (
	"errors"
	"reflect"
)

func ArrValueTranslate(from interface{}, to interface{}) (err error) {

	//reflect.ValueOf  获取数据源的数据
	//reflect.Indirect 假如这个值是一个指针  则取值，否则直接使用
	fromVal := reflect.Indirect(reflect.ValueOf(from))

	if fromVal.Kind() != reflect.Slice {
		err = errors.New("源数据不是数组")
	}

	if fromVal.Len() <= 0 {
		//源数据为空，不需要进行数据转移操作
		return
	}

	//获取输出对象的可变更指针
	toVal := reflect.ValueOf(to).Elem()
	//获取输出对象的类型
	toValType := toVal.Type()
	//获取输出对象存储数据的类型
	toValElemType := toValType.Elem()

	if toValElemType.Kind() != reflect.Interface && reflect.ValueOf(fromVal.Index(0).Interface()).Type() != toValElemType {
		err = errors.New("数据源和目标数据数据类型不一致")
		return
	}

	//用存储类型构建一个slice
	sliceType := reflect.SliceOf(toValElemType)
	valSlice := reflect.MakeSlice(sliceType, fromVal.Len(), fromVal.Len())
	//循环遍历出数据源的数据，并设定到刚创建的slice中
	for i := 0; i < fromVal.Len(); i++ {
		tempInterface := fromVal.Index(i)
		currentData := reflect.ValueOf(tempInterface.Interface())

		currentField := valSlice.Index(i)
		currentField.Set(currentData)
	}
	toVal.Set(valSlice)
	return
}
