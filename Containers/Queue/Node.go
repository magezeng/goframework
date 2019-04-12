package Queue

import "fmt"

type Node struct {
	key   string
	value interface{}
}

func (n *Node) GetKey() string {
	return n.key
}

func (n *Node) GetValue() interface{} {
	return n.value
}

func (n *Node) String() string {
	return fmt.Sprintf("结点 [%s : %v] ", n.GetKey(), n.GetValue())
}

func NewNode(key string, value interface{}) *Node {
	return &Node{key: key, value: value}
}
