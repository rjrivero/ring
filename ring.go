/*
Package ring provides support for circular fixed-size data structures,
like (LIFO) Stacks and (FIFO) Queues, that wrap around once they reach
a certain maximum size.

Rings should be paired with arrays or slices. The ring will manage
the ring head (the oldest item pushed to the ring) and tail
(the most recent item pushed to the ring) as offsets into the
array or slice, e.g:

  buffer := [8]string
  r := ring.New(8)

  // Push something into the buffer
  buffer[r.Push()] = "First item!"
  // Every time you push, the tail of the ring moves forward
  buffer[r.Push()] = "Second item!"
  buffer[r.Push()] = "Third item!"
  buffer[r.Push()] = "Fourth item!"
  buffer[r.Push()] = "Fifth item!"

  // r.Pop() will return the offset in the buffer
  // of the item most recently pushed
  print(buffer[r.Pop()])
  // Output: Fifth item!

  // If you keep popping, the tail of the ring
  // will keep going backwards.
  print(buffer[r.Pop()])
  // Output: Fourth item!

  // r.PopFront() will return the offset in the buffer
  // of the oldest item.
  print(buffer[r.PopFront()])
  // Output: First item!

  // If you keep front-popping, the head of the ring
  // will keep moving forward.
  print(buffer[r.PopFront()])
  // Output: Second item!

If you keep pushing items, the tail of the ring will eventually
wrap around the end of the slice, until it reaches the head.
Once that happens, the ring is full and pushing more items will
evict the oldest ones, moving the head of the ring forward.
*/
package ring

// Ring struct manages the indexes for a circular buffer.
type Ring struct {
	size      int
	tail, len int
}

// New initializes a Ring with the proper size.
func New(size int) Ring {
	return Ring{size: size}
}

// Push returns the current tail of the ring, and moves it forward
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

// Pop returns the previous tail of the ring, or -1 if empty.
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

// PopFront returns the head of the ring, or -1 if empty.
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

// Len returns the number of entries in the ring.
func (r Ring) Len() int {
	return r.len
}

// Cap returns the ring capacity (size).
func (r Ring) Cap() int {
	return r.size
}

// Full returns true if ring is full.
func (r Ring) Full() bool {
	return r.len == r.size
}

// Some returns true if the ring is not empty.
func (r Ring) Some() bool {
	return r.len > 0
}
