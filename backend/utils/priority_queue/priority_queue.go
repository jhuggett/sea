package priority_queue

import "container/heap"

type Compareable interface {
	SameAs(other Compareable) bool
}

// An Item is something we manage in a priority queue.
type Item[T Compareable] struct {
	value T // The value of the item.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T Compareable] struct {
	items   []*Item[T]
	compare func(a, b T) bool
}

func (pq PriorityQueue[T]) Len() int { return len(pq.items) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority.
	return pq.compare(pq.items[i].value, pq.items[j].value)
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	if len(pq.items) == 0 {
		return
	}

	pq.items[i], pq.items[j] = pq.items[j], pq.items[i]
	pq.items[i].index = i
	pq.items[j].index = j
}

// Don't use this method directly. Use PushIt instead. This is here to satisfy the heap.Interface which was pre-generics.
func (pq *PriorityQueue[T]) Push(x any) {
	item := x.(*Item[T])
	item.index = len(pq.items)
	pq.items = append(pq.items, item)
}

// Don't use this method directly. Use PopIt instead. This is here to satisfy the heap.Interface which was pre-generics.
func (pq *PriorityQueue[T]) Pop() any {
	if len(pq.items) == 0 {
		return nil
	}

	old := pq.items
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	pq.items = old[0 : n-1]

	return item
}

func New[T Compareable](compare func(a, b T) bool) *PriorityQueue[T] {
	pq := &PriorityQueue[T]{
		items:   []*Item[T]{},
		compare: compare,
	}

	return pq
}

// PushIt pushes a new item onto the priority queue.
func (pq *PriorityQueue[T]) PushIt(value T) {
	heap.Push(pq, &Item[T]{value: value})
}

// PopIt removes the highest priority item from the queue and returns it.
func (pq *PriorityQueue[T]) PopIt() *T {
	if pq.Len() == 0 {
		return nil
	}
	return &heap.Pop(pq).(*Item[T]).value
}

func (pq *PriorityQueue[T]) RemoveIt(value T) {
	for i, item := range pq.items {
		if item.value.SameAs(value) {
			heap.Remove(pq, i)
			return
		}
	}
}
