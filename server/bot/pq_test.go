package bot

import (
	"container/heap"
	"testing"
)

func TestPriorityQueue(t *testing.T) {
	q := &PriorityQueue{}
	heap.Init(q)
	heap.Push(q, &Item{
		State: &State{
			MyScore: 3,
		},
		Priority: 1,
	})
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
		Priority: 2,
	})
	item := heap.Pop(q).(*Item)
	if item.State.MyScore != 2 {
		t.Error(item.State.MyScore)
	}
}
