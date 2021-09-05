package bot

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	q := &PriorityQueue{
		arr: make([]*Item, 10_000),
	}
	heap.Init(q)
	heap.Push(q, &Item{
		State: &State{
			MyScore: 1,
		},
		Priority: 0,
	})
	heap.Push(q, &Item{
		State: &State{
			MyScore: 2,
		},
		Priority: 1,
	})
	item := heap.Pop(q).(*Item)
	if item.State.MyScore != 2 {
		t.Error(item.State.MyScore)
	}
}
