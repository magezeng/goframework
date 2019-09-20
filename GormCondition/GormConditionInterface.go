package GormCondition

type GormConditionInterface interface {
	Content() string
	Or(content string, args ...interface{}) GormConditionInterface
	And(content string, args ...interface{}) GormConditionInterface
	Group() GormConditionInterface
	IsNotNull() bool
	SetNextOption(option string)
}
