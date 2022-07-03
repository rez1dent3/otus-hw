package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("cache capacity=0", func(t *testing.T) {
		c := NewCache(0)
		c.Set("hello", "world")

		val, ok := c.Get("hello")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("popping an element", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		val, ok := c.Get("a")
		require.True(t, ok)
		require.Equal(t, 1, val)

		c.Set("b", 2)
		val, ok = c.Get("b")
		require.True(t, ok)
		require.Equal(t, 2, val)

		c.Set("c", 3)
		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 3, val)

		c.Set("d", 4) // drop a
		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 4, val)

		val, ok = c.Get("a")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("pareto rule", func(t *testing.T) {
		c := NewCache(4)
		c.Set("4", 4) // [4]
		c.Set("3", 3) // [3, 4]
		c.Set("2", 2) // [2, 3, 4]
		c.Set("1", 1) // [1, 2, 3, 4]

		c.Get("2") // [2, 1, 3, 4]
		c.Get("2") // [2, 1, 3, 4]
		c.Get("3") // [3, 2, 1, 4]

		c.Set("5", 5) // [5, 3, 2, 1]

		_, ok := c.Get("1")
		require.True(t, ok)

		_, ok = c.Get("2")
		require.True(t, ok)

		_, ok = c.Get("3")
		require.True(t, ok)

		_, ok = c.Get("4")
		require.False(t, ok) // not exists

		_, ok = c.Get("5")
		require.True(t, ok)
	})

	t.Run("purge logic", func(t *testing.T) {
		c := NewCache(3)

		c.Set("a", 1)
		c.Set("b", 2)
		c.Set("c", 3)

		c.Clear()

		val, ok := c.Get("b")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("c")
		require.False(t, ok)
		require.Nil(t, val)

		val, ok = c.Get("d")
		require.False(t, ok)
		require.Nil(t, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
