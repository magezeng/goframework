package GormCondition

type GormConditionGroup struct {
	condition  GormConditionInterface
	nextOption string
}

func (group *GormConditionGroup) SetNextOption(option string) {
	group.nextOption = option
}
func (group *GormConditionGroup) IsNotNull() bool {
	if group.condition == nil {
		return false
	}
	return group.condition.IsNotNull()
}

func (group *GormConditionGroup) Content() string {
	return "( " + group.condition.Content() + " )" + " " + group.nextOption + " "
}

func (group *GormConditionGroup) Or(content string, args ...interface{}) GormConditionInterface {
	group.SetNextOption("OR")
	return &GormConditionList{group, &GormCondition{content: content, args: args}}
}

func (group *GormConditionGroup) And(content string, args ...interface{}) GormConditionInterface {
	group.SetNextOption("AND")
	return &GormConditionList{group, &GormCondition{content: content, args: args}}
}
func (group *GormConditionGroup) Group() GormConditionInterface {
	return group
}
