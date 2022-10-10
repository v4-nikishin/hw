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

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		l.PushBack(20)     // [10, 30, 20]
		front := l.Front() // 10
		l.Remove(front)    // [30, 20]
		require.Equal(t, 2, l.Len())

		back := l.Back() // 20
		l.Remove(back)   // [30]
		require.Equal(t, 1, l.Len())

		back = l.Back() // nil
		l.Remove(back)  // []
		require.Equal(t, 0, l.Len())

		for i, v := range [...]int{80, 60, 40, 10, 30, 50, 70} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [70, 30, 40, 80, 60, 10, 50]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 70, l.Front().Value)
		require.Equal(t, 50, l.Back().Value)

		l.MoveToFront(l.Front()) // [70, 30, 40, 80, 60, 10, 50]
		l.MoveToFront(l.Back())  // [50, 70, 30, 40, 80, 60, 10]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 70, 30, 40, 80, 60, 10}, elems)
	})
}
