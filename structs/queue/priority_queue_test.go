package queue_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/structx/go-pkg/structs/queue"
)

func Test_Push(t *testing.T) {
	t.Run("push", func(_ *testing.T) {
		pq := queue.NewPriorityQueue()
		pq.Push(1, "hello")
	})
}

func Test_Pop(t *testing.T) {
	t.Run("pop", func(t *testing.T) {

		assert := assert.New(t)

		pq := queue.NewPriorityQueue()
		pq.Push(1, "hello")

		i := pq.Pop()
		assert.Equal(int64(1), i.GetPriority())
	})
}

func Test_Peek(t *testing.T) {
	t.Run("peek", func(t *testing.T) {

		assert := assert.New(t)

		pq := queue.NewPriorityQueue()
		pq.Push(1, "hello")

		i := pq.Peek()
		assert.Equal(int64(1), i.GetPriority())
	})
}
