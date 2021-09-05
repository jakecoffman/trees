package bot

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
}

func (pq PriorityQueue) Len() int {
	return len(pq.arr)
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
	n := len(pq.arr)
	item := x.(*Item)
	item.Index = n
	pq.arr = append(pq.arr, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	n := len(pq.arr)
	item := pq.arr[n-1]
	pq.arr[n-1] = nil // avoid memory leak
	item.Index = -1   // for safety
	pq.arr = pq.arr[0 : n-1]
	return item
}
