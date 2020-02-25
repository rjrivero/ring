package test

// ModRing is an alternate Ring implementation using MOD operator (%).
// This implementation has been discarded because of lower performance.
type ModRing struct {
	size      int
	tail, len int
}

// ModNew initializes a Ring with the proper size
func ModNew(size int) ModRing {
	return ModRing{size: size}
}

// Push returns a slot at the tail of the queue.
func (r *ModRing) Push() int {
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
func (r *ModRing) Pop() int {
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
func (r *ModRing) PopFront() int {
	if r.len == 0 {
		return -1
	}
	head := (r.size + r.tail - r.len) % r.size
	r.len--
	return head
}

// Len returns the current size of the ring.
func (r ModRing) Len() int {
	return r.len
}

// Cap returns the ring capacity (size)
func (r ModRing) Cap() int {
	return r.size
}

// Full returns true if ring is full
func (r ModRing) Full() bool {
	return r.len == r.size
}

// Some returns true if the ring is not empty
func (r ModRing) Some() bool {
	return r.len > 0
}
