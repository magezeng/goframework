package GormCondition

type Interface interface {
	Content() string
	Or(content string, args ...interface{}) Interface
	And(content string, args ...interface{}) Interface
	Group() Interface
	IsNotNull() bool
	SetNextOption(option string)
}

func NewGormCondition() Interface {
	return &GormCondition{}
}
