package main

// Item is a thing in the queue
type Item struct {
	// State is the value of the item; arbitrary.
	*State
	// Priority effects how it is sorted (min or max)
	Priority int
	// Index is needed by update and is maintained by the heap.Interface methods.
	Index int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue struct {
	arr []*Item
	cur int
}

func (pq PriorityQueue) Len() int {
	return pq.cur
}

func (pq PriorityQueue) Less(i, j int) bool {
	// higher priority first
	return pq.arr[i].Priority > pq.arr[j].Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq.arr[i], pq.arr[j] = pq.arr[j], pq.arr[i]
	pq.arr[i].Index = i
	pq.arr[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(*Item)
	item.Index = pq.cur
	pq.arr[pq.cur] = item
	pq.cur++
}

func (pq *PriorityQueue) Pop() interface{} {
	pq.cur--
	return pq.arr[pq.cur]
}
