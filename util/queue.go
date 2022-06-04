package util

import "container/list"

// IsContainsInt 包含数字
func IsContainInt(items []int, item int) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// Queue 队列
type Queue struct {
	list *list.List
}

// New  返回一个新队列
func New() *Queue {
	return &Queue{list.New()}
}

// PushBack adds element to the back of the queue
func (queue *Queue) PushBack(value interface{}) {
	queue.list.PushBack(value)
}

// PushFront adds element to the front of the queue
func (queue *Queue) PushFront(value interface{}) {
	queue.list.PushFront(value)
}

// Front returns first element of the queue
func (queue *Queue) Front() interface{} {
	it := queue.list.Front()
	if it != nil {
		return it.Value
	}
	return nil
}

func (queue *Queue) FrontRaw() *list.Element {
	return queue.list.Front()
}

// Back returns last element of the queue
func (queue *Queue) Back() interface{} {
	it := queue.list.Back()
	if it != nil {
		return it.Value
	}
	return nil
}

// Pop returns and deletes first element of the queue
func (queue *Queue) Pop() interface{} {
	it := queue.list.Front()
	if it != nil {
		queue.list.Remove(it)
		return it.Value
	}
	return nil
}

// Size returns size of the queue
func (queue *Queue) Size() int {
	return queue.list.Len()
}

// Empty returns whether queue is empty
func (queue *Queue) Empty() bool {
	return queue.list.Len() == 0
}

// Clear clears the queue
func (queue *Queue) Clear() {
	for !queue.Empty() {
		queue.Pop()
	}
}
