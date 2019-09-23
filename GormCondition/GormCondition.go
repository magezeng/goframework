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
	tempContent := condition.content
	for _, arg := range condition.args {
		tempContent = strings.Replace(tempContent, "=?=", fmt.Sprint(arg), 1)
	}
	return " " + tempContent + " " + condition.nextOption
}

func (condition *GormCondition) Or(content string, args ...interface{}) Interface {
	if condition.IsNotNull() {
		condition.SetNextOption("OR")
		return &List{condition, &GormCondition{content: content, args: args}}
	} else {
		return &GormCondition{content: content, args: args}
	}
}

func (condition *GormCondition) And(content string, args ...interface{}) Interface {
	if condition.IsNotNull() {
		condition.SetNextOption("AND")
		return &List{condition, &GormCondition{content: content, args: args}}
	} else {
		return &GormCondition{content: content, args: args}
	}
}

func (condition *GormCondition) Group() Interface {
	return &Group{condition: condition}
}
