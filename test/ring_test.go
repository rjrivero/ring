package test

import (
	"fmt"
	"testing"

	"github.com/rjrivero/ring"
)

// test the ring by pushing and popping alternatively
func ringTest(t *testing.T, ringSize int, pushPops []int) {
	ring := ring.NewRing(ringSize)
	head, tail, size := 0, 0, 0
	push := func(times int) {
		for i := 0; i < times; i++ {
			if pos := ring.Push(); pos != tail {
				t.Errorf("tail at repetition %d should be %d, got %d", i, tail, pos)
				return
			}
			tail = (tail + 1) % ringSize
			size++
			if size >= ringSize {
				// We are dragging head alongside
				size = ringSize
				head = tail
			}
		}
	}
	pop := func(times int) {
		for i := 0; i < times; i++ {
			if pos := ring.Pop(); pos != head {
				t.Errorf("head at repetition %d should be %d, got %d", i, head, pos)
				return
			}
			head = (head + 1) % ringSize
			size--
		}
	}
	// Alternate pushing and popping
	actions := []func(int){push, pop}
	for index, times := range pushPops {
		action := actions[index%2]
		action(times)
		if !t.Failed() {
			ringMetrics(t, ring, head, size, ringSize)
		}
		if t.Failed() {
			t.Logf("At push/pop %d", index)
			return
		}
	}
}

// Check ring metrics against expectations
func ringMetrics(t *testing.T, r ring.Ring, ringHead, ringLen, ringSize int) {
	if rLen := r.Len(); rLen != ringLen {
		t.Errorf("Length should be %d, got %d", ringLen, rLen)
	}
	if rCap := r.Cap(); rCap != ringSize {
		t.Errorf("Size should be %d, got %d", ringSize, rCap)
	}
	if full := r.Full(); ringLen == ringSize && !full {
		t.Errorf("Full should be true, got %t", full)
	}
	if full := r.Full(); ringLen != ringSize && full {
		t.Errorf("Full should be false, got %t", full)
	}
	if ringLen == 0 {
		if pos := r.Pop(); pos != -1 {
			t.Errorf("Empty Pop should yield -1, got %d", pos)
		}
	} else if head := r.Head(); head != ringHead {
		t.Errorf("Head should be %d, got %d", ringHead, head)
	}
	if t.Failed() {
		return
	}
	save := r
	iterMetrics(t, r.Iter(), ringHead, ringLen, ringSize)
	if r != save {
		t.Error("Iterator should not modify ring")
	}
}

// Check iterator metrics against expectations
func iterMetrics(t *testing.T, iter ring.Iterator, head, ringLen, ringSize int) {
	if ringLen <= 0 {
		if iter.Next() {
			t.Errorf("Iterator should be empty")
		}
		return
	}
	for i := 0; i < ringLen; i++ {
		if !iter.Next() {
			t.Error("Iterator should not be empty")
			return
		}
		expect := (head + i) % ringSize
		if pos := iter.Pos(); pos != expect {
			t.Errorf("Iteration number %d should yield %d, got %d", i, expect, pos)
			return
		}
	}
	if iter.Next() {
		t.Error("Iterator should be exhausted")
	}
}

func TestRing(t *testing.T) {
	type test struct {
		label    string
		pushPops []int
	}
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		wrap := (size * 3) / 2
		half := (size / 2) + 1
		tests := []test{
			{label: "Empty ring", pushPops: []int{0}},
			{label: "Full ring", pushPops: []int{size}},
			{label: "Wrap ring", pushPops: []int{wrap}},
			{label: "Deplete ring", pushPops: []int{size, size}},
			{label: "Pop some", pushPops: []int{size - 2, 2}},
			{label: "Invert", pushPops: []int{half, half - 1, half + 1}},
		}
		for _, current := range tests {
			t.Run(fmt.Sprintf("%s [size %d]", current.label, size), func(t *testing.T) {
				ringTest(t, size, current.pushPops)
			})
		}
	}
}
