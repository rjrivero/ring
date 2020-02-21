// Package ring provides support for circular queues backed
// by fixed-size buffers.
package ring

// Ring struct manages the indexes for a circular Ring queue.
type Ring struct {
	size       int
	head, tail int
	full       bool
}

// NewRing initializes a Ring with the proper size
func NewRing(size int) Ring {
	return Ring{size: size}
}

// Iterator supports iteration over a Ring
type Iterator struct {
	cursor, size, left int
}

// Push returns the tail of the queue and advances it
func (r *Ring) Push() int {
	pos := r.tail
	r.tail++
	if r.tail >= r.size {
		r.tail = 0
	}
	switch {
	case r.full:
		r.head = r.tail
	case r.head == r.tail:
		r.full = true
	}
	return pos
}

// Pop returns the head of the queue and advances it.
// returns -1 if the queue is empty.
func (r *Ring) Pop() int {
	switch {
	case r.full:
		r.full = false
	case r.head == r.tail:
		return -1
	}
	pos := r.head
	r.head++
	if r.head >= r.size {
		r.head = 0
	}
	return pos
}

// Len returns the current size of the ring.
func (r Ring) Len() int {
	switch {
	case r.full:
		return r.size
	case r.tail < r.head:
		return r.tail + r.size - r.head
	default:
		return r.tail - r.head
	}
}

// Cap returns the ring capacity (size)
func (r Ring) Cap() int {
	return r.size
}

// Full returns true if ring is full
func (r Ring) Full() bool {
	return r.full
}

// Iter builds an Iterator.
// The Ring should not be changed while iterated,
// otherwise results might be inconsistent.
func (r Ring) Iter() Iterator {
	return Iterator{
		size:   r.size,
		cursor: r.head - 1,
		left:   r.Len(),
	}
}

// Next returns true if the iterator is not exhausted
func (r *Iterator) Next() bool {
	if r.left <= 0 {
		return false
	}
	r.left--
	r.cursor++
	if r.cursor >= r.size {
		r.cursor = 0
	}
	return true
}

// Pos returns the current position in the queue.
func (r Iterator) Pos() int {
	return r.cursor
}
