package GormCondition

type List []Interface

func (list *List) SetNextOption(option string) {
	(*list)[len(*list)].SetNextOption(option)
}

func (list *List) IsNotNull() bool {
	return len(*list) > 0
}

func (list *List) Content() string {
	tempContent := ""
	for _, condition := range *list {
		tempContent += condition.Content()
	}
	return tempContent
}

func (list *List) Or(content string, args ...interface{}) Interface {
	(*list)[len(*list)].SetNextOption("OR")
	tempList := append(*list, &GormCondition{content: content, args: args})
	return &tempList
}

func (list *List) And(content string, args ...interface{}) Interface {
	(*list)[len(*list)].SetNextOption("AND")
	tempList := append(*list, &GormCondition{content: content, args: args})
	return &tempList
}

func (list *List) Group() Interface {
	return &Group{condition: list}
}
