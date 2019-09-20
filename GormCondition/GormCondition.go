package GormCondition

import (
	"fmt"
	"strings"
)

type GormCondition struct {
	content    string
	args       []interface{}
	nextOption string
}

func (condition *GormCondition) IsNotNull() bool {
	return len(condition.content) > 0
}

func (condition *GormCondition) SetNextOption(option string) {
	condition.nextOption = option
}

func (condition *GormCondition) Content() string {
	//	被填充内容用"=?="进行站位
	temp_content := condition.content
	for _, arg := range condition.args {
		temp_content = strings.Replace(temp_content, "=?=", fmt.Sprint(arg), 1)
	}
	return " " + temp_content + " " + condition.nextOption
}

func (condition *GormCondition) Or(content string, args ...interface{}) GormConditionInterface {
	if condition.IsNotNull() {
		condition.SetNextOption("OR")
		return &GormConditionList{condition, &GormCondition{content: content, args: args}}
	} else {
		return &GormCondition{content: content, args: args}
	}
}

func (condition *GormCondition) And(content string, args ...interface{}) GormConditionInterface {
	if condition.IsNotNull() {
		condition.SetNextOption("AND")
		return &GormConditionList{condition, &GormCondition{content: content, args: args}}
	} else {
		return &GormCondition{content: content, args: args}
	}
}

func (condition *GormCondition) Group() GormConditionInterface {
	return &GormConditionGroup{condition: condition}
}

func NewGormCondition() GormConditionInterface {
	return &GormCondition{}
}
