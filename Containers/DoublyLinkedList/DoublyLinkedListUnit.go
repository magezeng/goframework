package DoublyLinkedList

import (
	"errors"
)

type DoublyLinkedListUnit struct {
	element  *CanBeLinkedListItemInterface
	preUnit  *DoublyLinkedListUnit
	nextUnit *DoublyLinkedListUnit
}
type TestElement struct {
	Name string
	Age  float32
}

func (element *TestElement) GetIndex() interface{} {
	return element.Name
}

func (element *TestElement) GreatherWithIndex(compareIndex interface{}) bool {
	return element.GetIndex().(string) > compareIndex.(string)
}
func (unit DoublyLinkedListUnit) Next() (next *DoublyLinkedListUnit, err error) {
	if unit.nextUnit == nil {
		err = errors.New("没有下一项")
		return
	}
	next = unit.nextUnit
	return
}

func (unit *DoublyLinkedListUnit) Pre() (pre *DoublyLinkedListUnit, err error) {
	if unit.preUnit == nil {
		err = errors.New("没有下一项")
		return
	}
	pre = unit.preUnit
	return
}

func (unit *DoublyLinkedListUnit) GetIndex() (index interface{}) {
	index = (*unit.element).GetIndex()
	return
}

func (unit *DoublyLinkedListUnit) Update(newElement CanBeLinkedListItemInterface) (err error) {
	unit.element = &newElement
	return
}

func (unit *DoublyLinkedListUnit) GreatherWithIndex(compareIndex interface{}) bool {
	return (*unit.element).IsFrontThanIndex(compareIndex)
}

func (unit *DoublyLinkedListUnit) GreatherWithData(compareElement CanBeLinkedListItemInterface) bool {
	return (*unit.element).IsFrontThanIndex(compareElement.GetIndex())
}
