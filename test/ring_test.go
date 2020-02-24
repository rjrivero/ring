package test

import (
	"fmt"
	"testing"

	"github.com/rjrivero/ring"
)

// Ring interface for blackbox testing ring structs
type Ring interface {
	Push() int
	Pop() int
	PopFront() int
	Len() int
	Cap() int
	Full() bool
	Some() bool
	Iter() Ring
}

type testRing struct {
	ring.Ring
}

// Iter returns a copy of the ring for iteration
func (r *testRing) Iter() Ring {
	i := *r
	return &i
}

// Cursor represents a bounded-size segment inside an infinite buffer
type cursor struct {
	base, len, cap int
}

// modSize returns mod(value, cursor.cap)
func (c cursor) modCap(value int) int {
	for value < 0 {
		value += c.cap
	}
	return value % c.cap
}

// Head of the cursor
func (c cursor) head() int {
	return c.modCap(c.base)
}

// last member of the segment
func (c cursor) last() int {
	return c.modCap(c.base + c.len - 1)
}

// one past the last member of the segment
func (c cursor) tail() int {
	return c.modCap(c.base + c.len)
}

// Push to the ring the given number of times
func (c *cursor) Push(r Ring) error {
	if tail, pos := c.tail(), r.Push(); tail != pos {
		return fmt.Errorf("tail should be %d, got %d", tail, pos)
	}
	if c.len < c.cap {
		c.len++
	} else {
		c.base++
	}
	return nil
}

// Pop from the ring the given number of times
func (c *cursor) Pop(r Ring) error {
	if last, pos := c.last(), r.Pop(); last != pos {
		return fmt.Errorf("tail should be %d, got %d", last, pos)
	}
	c.len--
	return nil
}

// PopFront from the ring the given number of times
func (c *cursor) PopFront(r Ring) error {
	if head, pos := c.head(), r.PopFront(); head != pos {
		return fmt.Errorf("head should be %d, got %d", head, pos)
	}
	c.base++
	c.len--
	return nil
}

// test the ring by pushing and popping alternatively
func ringTest(t *testing.T, ringSize int, r Ring, pushPops []int) {
	c := cursor{cap: ringSize}
	// Alternate pushing, popping and poppingFront
	actions := []func(Ring) error{c.Push, c.Pop, c.PopFront}
	for index, times := range pushPops {
		var step int
		for step = 0; step < times; step++ {
			if err := actions[index%3](r); err != nil {
				t.Error(err)
				break
			}
		}
		if !t.Failed() {
			ringMetrics(t, r, c.head(), c.len, c.cap)
		}
		if t.Failed() {
			t.Errorf("At push/pop/popFront (%d, %d)", index, step)
			break
		}
	}
	save := r
	fifoMetrics(t, r.Iter(), c.head(), c.len, c.cap)
	lifoMetrics(t, r.Iter(), c.head(), c.len, c.cap)
	if r != save {
		t.Error("Ring should be copied by value")
	}
}

// Check ring metrics against expectations
func ringMetrics(t *testing.T, r Ring, ringHead, ringLen, ringSize int) {
	if rLen := r.Len(); rLen != ringLen {
		t.Errorf("Length should be %d, got %d", ringLen, rLen)
	}
	if rCap := r.Cap(); rCap != ringSize {
		t.Errorf("Size should be %d, got %d", ringSize, rCap)
	}
	var full, some bool
	switch {
	case ringLen == 0:
		full, some = false, false
		if pos := r.Pop(); pos != -1 {
			t.Errorf("Pop should be -1, got %d", pos)
		}
		if pos := r.PopFront(); pos != -1 {
			t.Errorf("PopFront should be -1, got %d", pos)
		}
	case ringLen == ringSize:
		full, some = true, true
	default:
		full, some = false, true
	}
	if r.Full() != full {
		t.Errorf("Full should be %t", full)
	}
	if r.Some() != some {
		t.Errorf("Some should be %t", some)
	}
	if t.Failed() {
		return
	}
}

// Check iterator metrics against expectations
func fifoMetrics(t *testing.T, iter Ring, ringHead, ringLen, ringSize int) {
	for i := 0; i < ringLen; i++ {
		if !iter.Some() {
			t.Errorf("Fifo should not be empty in step %d", i)
			return
		}
		if pos := iter.PopFront(); pos != ringHead {
			t.Errorf("Fifo should yield %d at step %d, got %d", ringHead, i, pos)
			return
		}
		ringHead = (ringHead + 1) % ringSize
	}
	if iter.Some() {
		t.Error("Fifo should be exhausted")
	}
}

// Check iterator metrics against expectations
func lifoMetrics(t *testing.T, iter Ring, ringHead, ringLen, ringSize int) {
	ringTail := (ringHead + ringLen - 1) % ringSize
	for i := 0; i < ringLen; i++ {
		if !iter.Some() {
			t.Errorf("Lifo should not be empty in step %d", i)
			return
		}
		if pos := iter.Pop(); pos != ringTail {
			t.Errorf("Lifo should yield %d at step %d, got %d", ringTail, i, pos)
			return
		}
		ringTail = (ringTail - 1 + ringSize) % ringSize
	}
	if iter.Some() {
		t.Error("Lifo should be exhausted")
	}
}

type testSequence struct {
	label    string
	pushPops []int
}

func generateTests(size int) []testSequence {
	wrap := (size * 3) / 2
	half := (size / 2) + 1
	return []testSequence{
		{label: "Empty ring", pushPops: []int{0}},
		{label: "Full ring", pushPops: []int{size}},
		{label: "Wrap ring", pushPops: []int{wrap}},
		{label: "Deplete stack", pushPops: []int{size, size}},
		{label: "Pop from stack", pushPops: []int{size - 2, 2}},
		{label: "Invert stack", pushPops: []int{half, half - 1, 0, half + 1}},
		{label: "Deplete queue", pushPops: []int{size, 0, size}},
		{label: "Pop from queue", pushPops: []int{size - 2, 0, 2}},
		{label: "Invert queue", pushPops: []int{half, 0, half - 1, half + 1}},
	}
}

func TestRing(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		for _, current := range generateTests(size) {
			t.Run(fmt.Sprintf("%s [size %d]", current.label, size), func(t *testing.T) {
				r := testRing{Ring: ring.New(size)}
				ringTest(t, size, &r, current.pushPops)
			})
		}
	}
}
