package priority_queue

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPriorityQueue(t *testing.T) {
	pq := New(func(a, b int) bool {
		return a < b
	})

	pq.PushIt(1)
	pq.PushIt(3)
	pq.PushIt(2)

	assert.Equal(t, 3, pq.Len())

	assert.Equal(t, 1, heap.Pop(pq).(*Item[int]).value)

	assert.Equal(t, 2, heap.Pop(pq).(*Item[int]).value)

	assert.Equal(t, 3, heap.Pop(pq).(*Item[int]).value)

	pq.PushIt(1)
	pq.PushIt(3)

	assert.Equal(t, 1, heap.Pop(pq).(*Item[int]).value)

	assert.Equal(t, 3, heap.Pop(pq).(*Item[int]).value)

	assert.Equal(t, 0, pq.Len())

	pq.PushIt(3)

	assert.Equal(t, 1, pq.Len())

	assert.Equal(t, 3, heap.Pop(pq).(*Item[int]).value)
}
