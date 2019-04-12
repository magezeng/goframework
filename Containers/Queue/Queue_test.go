package Queue

import (
	"testing"
)

var q = New()

func TestQueue_PushFront(t *testing.T) {
	q.PushFront(NewNode("1", "v1"))
	q.PushFront(NewNode("2", "v2"))
	q.PushFront(NewNode("3", "v3"))
	q.PushFront(NewNode("-1", "v-1"))
	q.PushFront(NewNode("-2", "v-2"))
	q.PushFront(NewNode("-3", "v-3"))
	t.Log("头插法: ", q.String())
}

func TestQueue_PopFront(t *testing.T) {
	q.PopFront()
	q.PopTail()
	q.PopFront()
	q.PopTail()
	q.PopFront()
	q.PopTail()
	t.Log("移除元素后", q.String())
}

func TestQueue_PushTail(t *testing.T) {
	q.PushTail(NewNode("1", "v1"))
	q.PushTail(NewNode("2", "v2"))
	q.PushTail(NewNode("3", "v3"))
	q.PushTail(NewNode("-1", "v-1"))
	q.PushTail(NewNode("-2", "v-2"))
	q.PushTail(NewNode("-3", "v-3"))
	t.Log("尾插法: ", q.String())
}

func TestQueue_PushMixed(t *testing.T) {
	q.PushTail(NewNode("-4", "v-4"))
	q.PushFront(NewNode("4", "v4"))
	t.Log("混合插入", q.String())
}

func TestQueue_Find(t *testing.T) {
	node1 := q.Find("-4", true)
	node2 := q.Find("-4", false)
	t.Log(node1.String())
	t.Log(node2.String())
}

func TestQueue_Overflow(t *testing.T){
	limitedQueue := NewWithLimit(3)
	limitedQueue.PushFront(NewNode("1", "v1"))
	limitedQueue.PushFront(NewNode("2", "v2"))
	limitedQueue.PushFront(NewNode("3", "v3"))
	limitedQueue.PushFront(NewNode("-1", "v-1"))
	limitedQueue.PushFront(NewNode("-2", "v-2"))
	limitedQueue.PushFront(NewNode("-3", "v-3"))
	t.Log(limitedQueue.String())
}
