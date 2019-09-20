package GormCondition

type GormConditionList []GormConditionInterface

func (list *GormConditionList) SetNextOption(option string) {
	(*list)[len(*list)].SetNextOption(option)
}
func (list *GormConditionList) IsNotNull() bool {
	return len(*list) > 0
}

func (list *GormConditionList) Content() string {
	tempContent := ""
	for _, condition := range *list {
		tempContent += condition.Content()
	}
	return tempContent
}

func (list *GormConditionList) Or(content string, args ...interface{}) GormConditionInterface {
	(*list)[len(*list)].SetNextOption("OR")
	tempList := GormConditionList(append(*list, &GormCondition{content: content, args: args}))
	return &tempList
}

func (list *GormConditionList) And(content string, args ...interface{}) GormConditionInterface {
	(*list)[len(*list)].SetNextOption("AND")
	tempList := GormConditionList(append(*list, &GormCondition{content: content, args: args}))
	return &tempList
}

func (list *GormConditionList) Group() GormConditionInterface {
	return &GormConditionGroup{condition: list}
}
