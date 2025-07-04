//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2025 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package priorityqueue

type supportedValueType interface {
	any | uint64
}

// Item represents a queue item supporting an optional additional Value
type Item[T supportedValueType] struct {
	ID       uint64
	Dist     float32
	Rescored bool
	Value    T
}

// Queue is a priority queue supporting generic item values
type Queue[T supportedValueType] struct {
	items []Item[T]
	less  func(items []Item[T], i, j int) bool
}

// NewMin constructs a priority queue which prioritizes items with smaller distance
func NewMin[T supportedValueType](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]Item[T], 0, capacity),
		less: func(items []Item[T], i, j int) bool {
			return items[i].Dist < items[j].Dist
		},
	}
}

// NewMin constructs a priority queue which prioritizes items with larger distance and smaller ID:
// - higher scores first
// - if tied, lower document id first
// Thus, the signs in the Less function are opposite for scores and ids
func NewMinWithId[T supportedValueType](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]Item[T], 0, capacity),
		less: func(items []Item[T], i, j int) bool {
			if items[i].Dist == items[j].Dist {
				return items[i].ID > items[j].ID
			}
			return items[i].Dist < items[j].Dist
		},
	}
}

// NewMax constructs a priority queue which prioritizes items with greater distance
func NewMax[T supportedValueType](capacity int) *Queue[T] {
	return &Queue[T]{
		items: make([]Item[T], 0, capacity),
		less: func(items []Item[T], i, j int) bool {
			return items[i].Dist > items[j].Dist
		},
	}
}

func (q *Queue[T]) ShouldEnqueue(distance float32, limit int) bool {
	return q.Len() < limit || q.Top().Dist < distance
}

func (q *Queue[T]) InsertAndPop(id uint64, score float64, limit int, worstDist *float64, val T) {
	q.InsertWithValue(id, float32(score), val)
	for q.Len() > limit {
		q.Pop()
	}
	// only update the worst distance when the queue is full, otherwise results can be missing if the first
	// entry that is checked already has a very high score
	if q.Len() >= limit {
		*worstDist = float64(q.Top().Dist)
	}
}

// Pop removes the next item in the queue and returns it
func (q *Queue[T]) Pop() Item[T] {
	out := q.items[0]
	q.items[0] = q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	q.heapify(0)
	return out
}

// Top peeks at the next item in the queue
func (q *Queue[T]) Top() Item[T] {
	return q.items[0]
}

// Len returns the length of the queue
func (q *Queue[T]) Len() int {
	return len(q.items)
}

// Cap returns the remaining capacity of the queue
func (q *Queue[T]) Cap() int {
	return cap(q.items)
}

// Reset clears all items from the queue
func (q *Queue[T]) Reset() {
	q.items = q.items[:0]
}

// ResetCap drops existing queue items, and allocates a new queue with the given capacity
func (q *Queue[T]) ResetCap(capacity int) {
	q.items = make([]Item[T], 0, capacity)
}

// Insert creates a valueless item and adds it to the queue
func (q *Queue[T]) Insert(id uint64, distance float32) int {
	item := Item[T]{
		ID:   id,
		Dist: distance,
	}
	return q.insert(item)
}

// InsertWithValue creates an item with a T type value and adds it to the queue
func (q *Queue[T]) InsertWithValue(id uint64, distance float32, val T) int {
	item := Item[T]{
		ID:    id,
		Dist:  distance,
		Value: val,
	}
	return q.insert(item)
}

// DeleteItem deletes item meeting predicate's conditions
// TODO aliszka optimize?
func (q *Queue[T]) DeleteItem(match func(item Item[T]) bool) bool {
	for i := range q.items {
		if match(q.items[i]) {
			if i == 0 {
				q.Pop()
			} else {
				last := q.Len() - 1
				q.items[i] = q.items[last]
				q.items = q.items[:last]
				q.heapify(0)
			}
			return true
		}
	}
	return false
}

func (q *Queue[T]) insert(item Item[T]) int {
	q.items = append(q.items, item)
	i := len(q.items) - 1
	for i != 0 && q.less(q.items, i, q.parent(i)) {
		q.swap(i, q.parent(i))
		i = q.parent(i)
	}
	return i
}

func (q *Queue[T]) left(i int) int {
	return 2*i + 1
}

func (q *Queue[T]) right(i int) int {
	return 2*i + 2
}

func (q *Queue[T]) parent(i int) int {
	return (i - 1) / 2
}

func (q *Queue[T]) swap(i, j int) {
	q.items[i], q.items[j] = q.items[j], q.items[i]
}

func (q *Queue[T]) heapify(i int) {
	left := q.left(i)
	right := q.right(i)
	smallest := i
	if left < len(q.items) && q.less(q.items, left, i) {
		smallest = left
	}

	if right < len(q.items) && q.less(q.items, right, smallest) {
		smallest = right
	}

	if smallest != i {
		q.swap(i, smallest)
		q.heapify(smallest)
	}
}
