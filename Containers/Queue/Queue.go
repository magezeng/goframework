package Queue

import (
	"bytes"
	"fmt"
)

// Queue 双端队列，节点数和长度不一定相等，用于动态增长或减少
type Queue struct {
	nodes []*Node
	// 记录头部下标
	front int
	// 记录尾部下标
	tail int
	// 记录长度
	length int
	// 最大长度
	maxLength int
}

func New() *Queue {
	return new(Queue).init(-1)
}

func NewWithLimit(maxLength int) *Queue {
	return new(Queue).init(maxLength)
}

// init 初始化方法，返回一个立刻可以使用的队列
func (q *Queue) init(maxLength int) *Queue {
	q.nodes = make([]*Node, 1)
	q.front, q.tail, q.length = 0, 0, 0
	q.maxLength = maxLength
	return q
}

// lazyInit 懒初始化
func (q *Queue) lazyInit() {
	if q.nodes == nil {
		q.init(-1)
	}
}

// Len 返回队列的长度
func (q *Queue) Len() int {
	return q.length
}

// IsEmpty 队列是否为空队列
func (q *Queue) IsEmpty() bool {
	return q.length == 0
}

// IsFull 队列是否为满队列
func (q *Queue) IsFull() bool {
	return q.length == len(q.nodes)
}

// sparse 判断是不是内存有浪费
func (q *Queue) sparse() bool {
	// 当长度远大于结点个数四倍以上时
	return 1 < q.length && q.length < len(q.nodes)/4
}

// resize 重新调整大小
func (q *Queue) resize(size int) {
	adjusted := make([]*Node, size)
	if q.front < q.tail {
		copy(adjusted, q.nodes[q.front:q.tail])
	} else {
		n := copy(adjusted, q.nodes[q.front:])
		copy(adjusted[n:], q.nodes[:q.tail])
	}
	q.nodes = adjusted
	q.front = 0
	q.tail = q.length
}

// lazyGrow 2倍增长
func (q *Queue) lazyGrow() {
	if q.IsFull() {
		q.resize(len(q.nodes) * 2)
	}
}

// lazyShrink 2倍压缩
func (q *Queue) lazyShrink() {
	if q.sparse() {
		q.resize(len(q.nodes) / 2)
	}
}

func (q *Queue) String() string {
	var result bytes.Buffer
	result.WriteByte('[')
	j := q.front
	for i := 0; i < q.length; i++ {
		result.WriteString(fmt.Sprintf("%v", q.nodes[j]))
		if i < q.length-1 {
			result.WriteByte(' ')
		}
		j = q.nextIndex(j)
	}
	result.WriteByte(']')
	return result.String()
}

// nextIndex 返回下一个下标
func (q *Queue) nextIndex(i int) int {
	return (i + 1) & (len(q.nodes) - 1)
}

// prevIndex 返回上一个下标
func (q *Queue) prevIndex(i int) int {
	return (i - 1) & (len(q.nodes) - 1)
}

// isLimited 是否有长度限制
func (q *Queue) isLimited() bool {
	return !(q.maxLength == -1)
}

// Front 返回队列头节点
func (q *Queue) Front() *Node {
	return q.nodes[q.front]
}

// Tail 返回队列尾部节点
func (q *Queue) Tail() *Node {
	return q.nodes[q.nextIndex(q.tail)]
}

// PushFront 头部插入结点
func (q *Queue) PushFront(v *Node) {
	q.lazyInit()
	q.lazyGrow()
	q.front = q.prevIndex(q.front)
	q.nodes[q.front] = v
	q.length++
	if q.isLimited() && q.length > q.maxLength {
		q.PopTail()
	}
}

// PushBack 尾部插入结点
func (q *Queue) PushTail(v *Node) {
	q.lazyInit()
	q.lazyGrow()
	q.nodes[q.tail] = v
	q.tail = q.nextIndex(q.tail)
	q.length++
	if q.isLimited() && q.length > q.maxLength {
		q.PopFront()
	}
}

// PopFront 从头部删除一个结点并返回这个结点
func (q *Queue) PopFront() *Node {
	if q.IsEmpty() {
		return nil
	}
	v := q.nodes[q.front]
	q.nodes[q.front] = nil
	q.front = q.nextIndex(q.front)
	q.length--
	q.lazyShrink()
	return v
}

// PopTail 从尾部删除一个结点并返回这个结点
func (q *Queue) PopTail() *Node {
	if q.IsEmpty() {
		return nil
	}
	q.tail = q.prevIndex(q.tail)
	v := q.nodes[q.tail]
	q.nodes[q.tail] = nil
	q.length--
	q.lazyShrink()
	return v
}

// Find 寻找对应键值的结点
func (q *Queue) Find(key string, fromFront bool) *Node {
	nodes := make([]*Node, len(q.nodes))
	copy(nodes, q.nodes)
	if !fromFront {
		for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
			nodes[i], nodes[j] = nodes[j], nodes[i]
		}
	}
	for _, node := range nodes {
		if node != nil && node.GetKey() == key {
			return node
		}
	}
	return nil
}
