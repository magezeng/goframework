package DoublyLinkedList

import (
	"errors"
	"sync"
)

type DoublyLinkedList struct {
	startUnit *DoublyLinkedListUnit
	endUnit   *DoublyLinkedListUnit
	indexMap  map[interface{}]*DoublyLinkedListUnit
	rwMutex   sync.RWMutex
}

func New() (newList *DoublyLinkedList) {
	newList = new(DoublyLinkedList)
	newList.indexMap = map[interface{}]*DoublyLinkedListUnit{}
	return
}

func (list *DoublyLinkedList) InsertOrUpdateDatas(elements []CanBeLinkedListItemInterface) (err error) {
	for _, tempElement := range elements {
		err = list.InsertOrUpdateData(tempElement)
		if err != nil {
			return
		}
	}
	return
}

func (list *DoublyLinkedList) InsertOrUpdateData(element CanBeLinkedListItemInterface) (err error) {

	if tempUnit, ok := list.indexMap[element.GetIndex()]; ok {
		tempUnit.Update(element)
	} else {
		err = list.insertData(element)
	}
	return
}

func (list *DoublyLinkedList) insertData(element CanBeLinkedListItemInterface) (err error) {

	if _, ok := list.indexMap[element.GetIndex()]; ok {
		err = errors.New("插入失败,之前已经存在该Index")
		return
	}

	newUnit := new(DoublyLinkedListUnit)
	newUnit.element = &element
	if len(list.indexMap) <= 0 {
		list.startUnit = newUnit
		list.endUnit = newUnit
	} else {
		var willPreUnit *DoublyLinkedListUnit = nil
		for index, unit := range list.indexMap {
			if !newUnit.GreatherWithIndex(index) && (willPreUnit == nil || willPreUnit.GreatherWithIndex(index)) {
				willPreUnit = unit
			}
		}
		if willPreUnit == nil {
			//这种情况代表新元素应该排在第一个，因为没有老数据比他应有的排序更靠前
			oldFirstUnit := list.startUnit
			newUnit.nextUnit = oldFirstUnit
			list.startUnit = newUnit
			oldFirstUnit.preUnit = newUnit
		} else {
			if willPreUnit == list.endUnit {
				//这种情况代表新元素应该插入到末尾
				willPreUnit.nextUnit = newUnit
				newUnit.preUnit = willPreUnit
				list.endUnit = newUnit
			} else {
				//新元素插入在中间
				willNextUnit := willPreUnit.nextUnit
				willPreUnit.nextUnit = newUnit
				newUnit.preUnit = willPreUnit
				newUnit.nextUnit = willNextUnit
				willNextUnit.preUnit = newUnit
			}
		}
	}
	list.indexMap[element.GetIndex()] = newUnit
	return
}

func (list *DoublyLinkedList) GetInverseOrderData(count int) (elements []CanBeLinkedListItemInterface, err error) {

	currentIndex := 0
	tempElements := make([]CanBeLinkedListItemInterface, count)

	for currentUnit := list.endUnit; currentUnit != nil && currentIndex < count; currentUnit = currentUnit.preUnit {
		tempElements[currentIndex] = *currentUnit.element
		currentIndex += 1
	}

	elements = tempElements[:currentIndex]
	return
}

func (list *DoublyLinkedList) GetOrderData(count int) (elements []CanBeLinkedListItemInterface, err error) {

	currentIndex := 0
	tempElements := make([]CanBeLinkedListItemInterface, count)

	for currentUnit := list.startUnit; currentUnit != nil && currentIndex < count; currentUnit = currentUnit.nextUnit {
		tempElements[currentIndex] = *currentUnit.element
		currentIndex += 1
	}

	elements = tempElements[:currentIndex]
	return
}

func (list *DoublyLinkedList) GetInverseOrderDataFrom(fromIndex interface{}, count int) (elements []CanBeLinkedListItemInterface, err error) {

	currentIndex := 0
	tempElements := make([]CanBeLinkedListItemInterface, count)
	var currentUnit *DoublyLinkedListUnit

	if tempUnit, have := list.indexMap[fromIndex]; have {
		currentUnit = tempUnit.preUnit
	} else {
		currentUnit = list.endUnit
	}
	for ; currentUnit != nil && currentIndex < count; currentUnit = currentUnit.preUnit {
		tempElements[currentIndex] = *currentUnit.element
		currentIndex += 1
	}

	elements = tempElements[:currentIndex]
	return
}

func (list *DoublyLinkedList) GetOrderDataFrom(fromIndex interface{}, count int) (elements []CanBeLinkedListItemInterface, err error) {

	currentIndex := 0
	tempElements := make([]CanBeLinkedListItemInterface, count)

	currentUnit := list.startUnit

	if tempUnit, have := list.indexMap[fromIndex]; have {
		currentUnit = tempUnit.nextUnit
	}

	for ; currentUnit != nil && currentIndex < count; currentUnit = currentUnit.nextUnit {
		tempElements[currentIndex] = *currentUnit.element
		currentIndex += 1
	}

	elements = tempElements[:currentIndex]
	return
}

func (list *DoublyLinkedList) GetInverseOrderAllData() (elements []CanBeLinkedListItemInterface, err error) {
	elements, err = list.GetInverseOrderData(list.GetCount())
	return
}

func (list *DoublyLinkedList) GetOrderAllData() (elements []CanBeLinkedListItemInterface, err error) {
	elements, err = list.GetOrderData(list.GetCount())
	return
}

func (list *DoublyLinkedList) GetCount() (count int) {
	count = len(list.indexMap)
	return
}

func (list *DoublyLinkedList) PopFirst() (element CanBeLinkedListItemInterface) {
	if list.startUnit == nil {
		return
	}

	firstUnit := list.startUnit
	element = *firstUnit.element
	list.startUnit = firstUnit.nextUnit
	list.startUnit.preUnit = nil
	delete(list.indexMap, (*firstUnit.element).GetIndex())
	return
}

func (list *DoublyLinkedList) PopLast() (element CanBeLinkedListItemInterface) {
	if list.endUnit == nil {
		return
	}

	lastUnit := list.endUnit
	element = *lastUnit.element
	list.endUnit = lastUnit.preUnit
	list.endUnit.nextUnit = nil
	delete(list.indexMap, (*lastUnit.element).GetIndex())
	return
}

func (list *DoublyLinkedList) DeleteDataWithIndex(index interface{}) {
	willDeleteData := list.indexMap[index]
	delete(list.indexMap, index)

	if willDeleteData.preUnit == nil {
		list.startUnit = willDeleteData.nextUnit
	} else {
		willDeleteData.preUnit.nextUnit = willDeleteData.nextUnit
	}

	if willDeleteData.nextUnit == nil {
		list.endUnit = willDeleteData.preUnit
	} else {
		willDeleteData.nextUnit.preUnit = willDeleteData.preUnit
	}
}

func (list *DoublyLinkedList) DeleteDataWithData(element CanBeLinkedListItemInterface) {
	list.DeleteDataWithIndex(element.GetIndex())
}
