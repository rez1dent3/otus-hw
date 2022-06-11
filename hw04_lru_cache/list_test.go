package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("push back", func(t *testing.T) {
		l := NewList()
		l.PushBack(1) // [1]
		l.PushBack(2) // [1, 2]
		l.PushBack(3) // [1, 2, 3]

		require.Equal(t, 3, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 3, l.Back().Value)
	})

	t.Run("push front", func(t *testing.T) {
		l := NewList()
		l.PushFront(3) // [3]
		l.PushFront(2) // [2, 3]
		l.PushFront(1) // [1, 2, 3]

		require.Equal(t, 3, l.Len())
		require.Equal(t, 1, l.Front().Value)
		require.Equal(t, 3, l.Back().Value)
	})

	t.Run("clean up to zero", func(t *testing.T) {
		l := NewList()
		el := l.PushFront(3) // [3]
		require.Equal(t, 1, l.Len())
		require.Equal(t, 3, l.Front().Value)
		require.Equal(t, 3, l.Back().Value)

		l.Remove(el) // []

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("check the boundary values", func(t *testing.T) {
		l := NewList()
		el3 := l.PushFront(3) // [3]
		el2 := l.PushFront(2) // [2, 3]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 3, l.Back().Value)

		l.Remove(el3) // [2]

		require.Equal(t, 1, l.Len())
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 2, l.Back().Value)

		el4 := l.PushBack(4) // [2, 4]

		require.Equal(t, 2, l.Len())
		require.Equal(t, 2, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)

		l.Remove(el2) // [4]

		require.Equal(t, 1, l.Len())
		require.Equal(t, 4, l.Front().Value)
		require.Equal(t, 4, l.Back().Value)

		l.Remove(el4) // []

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		l.PushBack(40)  // [10, 20, 30, 40]
		l.PushFront(0)  // [0, 10, 20, 30, 40]
		require.Equal(t, 5, l.Len())

		front := l.Front()
		require.Equal(t, 0, front.Value)

		l.Remove(front)
		require.Equal(t, 4, l.Len())
		require.Equal(t, 10, l.Front().Value)

		back := l.Back()
		require.Equal(t, 40, back.Value)

		l.Remove(back)
		require.Equal(t, 3, l.Len())
		require.Equal(t, 30, l.Back().Value)

		middle := l.Front().Next // 20
		require.Equal(t, 20, middle.Value)

		l.Remove(middle) // [10, 30]
		require.Equal(t, 2, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Back().Value)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Back()) // [70, 80, 60, 40, 10, 30, 50]
		require.Equal(t, 70, l.Front().Value)
		require.Equal(t, 50, l.Back().Value)

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
