// Package ring provides support for circular fixed-size data structures.
package ring

// Ring struct manages the indexes for a circular buffer.
type Ring struct {
	size      int
	tail, len int
}

// New initializes a Ring with the proper size
func New(size int) Ring {
	return Ring{size: size}
}

// Push returns a slot at the tail of the queue.
func (r *Ring) Push() int {
	if r.len < r.size {
		r.len++
	}
	tail := r.tail
	if r.tail++; r.tail == r.size {
		r.tail = 0
	}
	return tail
}

// Pop returns the tail of the ring, or -1 if empty.
func (r *Ring) Pop() int {
	if r.len == 0 {
		return -1
	}
	if r.tail--; r.tail < 0 {
		r.tail = r.size - 1
	}
	r.len--
	return r.tail
}

// PopFront returns the head of the ring, or -1 if empty
func (r *Ring) PopFront() int {
	if r.len == 0 {
		return -1
	}
	head := r.tail - r.len
	if head < 0 {
		head += r.size
	}
	r.len--
	return head
}

// Len returns the current size of the ring.
func (r Ring) Len() int {
	return r.len
}

// Cap returns the ring capacity (size)
func (r Ring) Cap() int {
	return r.size
}

// Full returns true if ring is full
func (r Ring) Full() bool {
	return r.len == r.size
}

// Some returns true if the ring is not empty
func (r Ring) Some() bool {
	return r.len > 0
}
