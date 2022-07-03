package hw04lrucache

import (
	"sync"
)

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
	mx    sync.RWMutex
	front *ListItem
	back  *ListItem
	len   int
}

func (l *list) Len() int {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.len
}

func (l *list) Front() *ListItem {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.front
}

func (l *list) Back() *ListItem {
	l.mx.RLock()
	defer l.mx.RUnlock()

	return l.back
}

func (l *list) unsafePushFront(v interface{}) *ListItem {
	l.len++
	item := l.front
	value := ListItem{Value: v, Next: item}
	if item != nil {
		item.Prev = &value
	}

	l.front = &value
	if l.back == nil {
		l.back = l.front
	}

	return l.front
}

func (l *list) PushFront(v interface{}) *ListItem {
	l.mx.Lock()
	defer l.mx.Unlock()

	return l.unsafePushFront(v)
}

func (l *list) PushBack(v interface{}) *ListItem {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.len++
	item := l.back
	value := ListItem{Value: v, Prev: item}
	if item != nil {
		item.Next = &value
	}

	l.back = &value
	if l.front == nil {
		l.front = l.back
	}

	return l.back
}

func (l *list) unsafeRemove(i *ListItem) {
	l.len--
	prev, next := i.Prev, i.Next
	if prev != nil && next != nil {
		next.Prev = prev
		prev.Next = next
		return
	}

	if next != nil {
		l.front = next
		l.front.Prev = nil
	}

	if prev != nil {
		l.back = prev
		l.back.Next = nil
	}

	if l.back != nil && l.back.Value == i.Value {
		l.back = nil
	}

	if l.front != nil && l.front.Value == i.Value {
		l.front = nil
	}
}

func (l *list) Remove(i *ListItem) {
	l.mx.Lock()
	defer l.mx.Unlock()

	l.unsafeRemove(i)
}

func (l *list) MoveToFront(i *ListItem) {
	l.mx.Lock()
	defer l.mx.Unlock()

	if i.Prev == nil {
		return
	}

	l.unsafeRemove(i)
	l.unsafePushFront(i.Value)
}

func NewList() List {
	return new(list)
}
