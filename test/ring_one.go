package test

// OneRing is an alternate Ring implementation very much
// like ZeroRing, but uses head = -1 to identify empty ring,
// instead of full ring.
type OneRing struct {
	size       int
	head, tail int
}

// OneNew initializes a OneRing with the proper size
func OneNew(size int) OneRing {
	return OneRing{size: size, head: -1}
}

// Push returns a slot at the tail of the queue.
func (r *OneRing) Push() int {
	pos := r.tail
	if r.tail++; r.tail >= r.size {
		r.tail = 0
	}
	switch {
	case r.head == pos: // was full
		r.head = r.tail
	case r.head < 0: // was empty
		r.head = pos
	}
	return pos
}

// Pop returns the tail of the ring, or -1 if empty.
func (r *OneRing) Pop() int {
	if r.head < 0 {
		return -1
	}
	if r.tail--; r.tail < 0 {
		r.tail = r.size - 1
	}
	if r.tail == r.head {
		r.head = -1
	}
	return r.tail
}

// PopFront returns the head of the ring, or -1 if empty
func (r *OneRing) PopFront() int {
	if r.head < 0 {
		return -1
	}
	pos := r.head
	if r.head++; r.head >= r.size {
		r.head = 0
	}
	if r.head == r.tail {
		r.head = -1
	}
	return pos
}

// Len returns the current size of the ring.
func (r OneRing) Len() int {
	switch {
	case r.head < 0:
		return 0
	case r.tail <= r.head:
		return r.tail + r.size - r.head
	default:
		return r.tail - r.head
	}
}

// Cap returns the ring capacity (size)
func (r OneRing) Cap() int {
	return r.size
}

// Full returns true if ring is full
func (r OneRing) Full() bool {
	return r.head == r.tail
}

// Some returns true if the ring is not empty
func (r OneRing) Some() bool {
	return r.head >= 0
}
