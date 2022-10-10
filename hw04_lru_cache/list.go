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
	len  int
	head *ListItem
	back *ListItem
}

func (l *list) Len() int {
	return l.len
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.back
}

func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.head == nil {
		l.head = item
		l.back = item
	} else {
		item.Next = l.head
		l.head.Prev = item
		l.head = item
	}
	l.len++
	return item
}

func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{
		Value: v,
	}
	if l.back == nil {
		l.head = item
		l.back = item
	} else {
		item.Prev = l.back
		l.back.Next = item
		l.back = item
	}
	l.len++
	return item
}

func (l *list) Remove(i *ListItem) {
	defer func() {
		if l.len > 0 {
			l.len--
		}
	}()
	if l.head == nil {
		return
	}
	if l.len == 1 {
		l.head = nil
		l.back = nil
		return
	}
	if i == l.Front() {
		i.Next.Prev = nil
		l.head = i.Next
		return
	}
	if i == l.Back() {
		i.Prev.Next = nil
		l.back = i.Prev
		return
	}
	i.Next.Prev = i.Prev
	i.Prev.Next = i.Next
}

func (l *list) MoveToFront(i *ListItem) {
	if i == nil {
		return
	}
	if l.head == nil {
		return
	}
	if l.len == 1 {
		return
	}
	if i == l.Front() {
		return
	}
	if i == l.Back() {
		i.Prev.Next = nil
		l.back = i.Prev
	} else {
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}
	i.Prev = nil
	i.Next = l.Front()
	l.Front().Prev = i
	l.head = i
}

func NewList() List {
	return new(list)
}
