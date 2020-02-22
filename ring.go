// Package ring provides support for circular fixed-size data structures
package ring

// Ring struct manages the indexes for a circular buffer
type Ring struct {
	size       int
	head, tail int
}

// New initializes a Ring with the proper size
func New(size int) Ring {
	return Ring{size: size}
}

// Push returns a slot at the tail of the queue.
func (r *Ring) Push() int {
	pos := r.tail
	if r.tail++; r.tail >= r.size {
		r.tail = 0
	}
	if r.tail == r.head {
		// The queue became full. As long as it remains full, the head
		// and the tail will point to the same position.
		// However, r.head == r.tail is the empty condition,
		// so we need some way to tell the queue is full. We settled on:
		// - Empty queue: r.head == r.tail
		// - Full queue: r.head = -1, r.tail doubles as both head and tail.
		r.head = -1
	}
	return pos
}

// Pop returns the tail of the ring, or -1 if empty.
func (r *Ring) Pop() int {
	switch {
	case r.head < 0:
		r.head = r.tail
	case r.head == r.tail:
		return -1
	}
	if r.tail--; r.tail < 0 {
		r.tail = r.size - 1
	}
	return r.tail
}

// PopFront returns the head of the ring, or -1 if empty
func (r *Ring) PopFront() int {
	switch {
	case r.head < 0: // full
		r.head = r.tail
	case r.head == r.tail: // empty
		return -1
	}
	pos := r.head
	if r.head++; r.head >= r.size {
		r.head = 0
	}
	return pos
}

// Len returns the current size of the ring.
func (r Ring) Len() int {
	switch {
	case r.head < 0:
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
	return r.head < 0
}

// Some returns true if the ring is not empty
func (r Ring) Some() bool {
	return r.head != r.tail
}
