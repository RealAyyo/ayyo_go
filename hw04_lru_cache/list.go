package hw04lrucache

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	len  int
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.len++

	if l.head == nil {
		node := &ListItem{Value: v}

		l.head = node
		l.tail = node
		return node
	}
	node := &ListItem{Value: v, Next: l.head}

	l.head.Prev = node
	l.head = node
	return node
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.len++

	if l.head == nil {
		node := &ListItem{Value: v, Next: nil, Prev: nil}

		l.head = node
		l.tail = node
		return node
	}
	node := &ListItem{Value: v, Next: nil, Prev: l.tail}

	l.tail.Next = node
	l.tail = node
	return node
}

func (l *list) Remove(i *ListItem) {
	if i == nil {
		return
	}

	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		l.head = i.Next
	}

	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		l.tail = i.Prev
	}

	l.len--
}

func (l *list) MoveToFront(i *ListItem) {
	l.Remove(i)

	i.Prev = nil
	i.Next = l.head
	l.head.Prev = i
	l.head = i

	l.len++
}

func NewList() List {
	return new(list)
}
