package GormCondition

type Group struct {
	condition  Interface
	nextOption string
}

func (group *Group) SetNextOption(option string) {
	group.nextOption = option
}

func (group *Group) IsNotNull() bool {
	if group.condition == nil {
		return false
	}
	return group.condition.IsNotNull()
}

func (group *Group) Content() string {
	return "( " + group.condition.Content() + " )" + " " + group.nextOption + " "
}

func (group *Group) Or(content string, args ...interface{}) Interface {
	group.SetNextOption("OR")
	return &List{group, &GormCondition{content: content, args: args}}
}

func (group *Group) And(content string, args ...interface{}) Interface {
	group.SetNextOption("AND")
	return &List{group, &GormCondition{content: content, args: args}}
}

func (group *Group) Group() Interface {
	return group
}
