package test

import (
	"fmt"
	"testing"

	"github.com/rjrivero/ring"
)

func TestEmpty(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Empty ring of size %d", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			if ring.Len() != 0 {
				t.Errorf("Length should be 0")
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			if ring.Full() {
				t.Error("Full should be false")
			}
			if pos := ring.Pop(); pos != -1 {
				t.Errorf("Empty Pop should yield -1, got %d", pos)
			}
			save := ring
			iter := ring.Iter()
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}

func TestFull(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Full ring of size %d", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			for i := 0; i < size; i++ {
				if pos := ring.Push(); pos != i {
					t.Errorf("Index after %d pushes should be %d, got %d", i+1, i, pos)
					t.FailNow()
				}
			}
			if ring.Len() != size {
				t.Errorf("Length should be %d", size)
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			if !ring.Full() {
				t.Error("Full should be true")
			}
			save := ring
			iter := ring.Iter()
			for i := 0; i < size; i++ {
				if !iter.Next() {
					t.Error("Iterator should not be empty")
					t.FailNow()
				}
				if pos := iter.Pos(); pos != i {
					t.Errorf("Iteration number %d should yield %d, got %d", i, i, pos)
					t.FailNow()
				}
			}
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}

func TestWrap(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Wrap ring of size %d around", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			wrap := size * 3 / 2
			for i := 0; i < wrap; i++ {
				if pos := ring.Push(); pos != i%size {
					t.Errorf("Index after %d pushes should be %d, got %d", i+1, i%size, pos)
					t.FailNow()
				}
			}
			if ring.Len() != size {
				t.Errorf("Length should be %d", size)
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			if !ring.Full() {
				t.Error("Full should be true")
			}
			save := ring
			iter := ring.Iter()
			wrap = wrap % size
			for i := 0; i < size; i++ {
				if !iter.Next() {
					t.Error("Iterator should not be empty")
					t.FailNow()
				}
				if pos := iter.Pos(); pos != wrap%size {
					t.Errorf("Iteration number %d should yield %d, got %d", i, wrap, pos)
					t.FailNow()
				}
				wrap++
			}
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}

func TestDeplete(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Deplete ring of size %d", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			wrap := size * 3 / 2
			for i := 0; i < wrap; i++ {
				if pos := ring.Push(); pos != i%size {
					t.Errorf("Index after %d pushes should be %d, got %d", i+1, i%size, pos)
					t.FailNow()
				}
			}
			wrap = wrap % size
			for i := 0; i < size; i++ {
				if pos := ring.Pop(); pos != wrap%size {
					t.Errorf("Index after %d pops should be %d, got %d", i+1, wrap%size, pos)
					t.FailNow()
				}
				wrap++
			}
			if pos := ring.Pop(); pos != -1 {
				t.Errorf("Pop once depleted should be -1, got %d", pos)
			}
			if ring.Len() != 0 {
				t.Errorf("Length should be %d", size)
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			if ring.Full() {
				t.Error("Full should be false")
			}
			save := ring
			iter := ring.Iter()
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}

func TestPopSome(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Pop some items from ring of size %d", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			gapTail := 2
			for i := 0; i < size-gapTail; i++ {
				if pos := ring.Push(); pos != i%size {
					t.Errorf("Index after %d pushes should be %d, got %d", i+1, i%size, pos)
					t.FailNow()
				}
			}
			gapHead := 2
			for i := 0; i < gapHead; i++ {
				if pos := ring.Pop(); pos != i {
					t.Errorf("Index after %d pops should be %d, got %d", i+1, i, pos)
					t.FailNow()
				}
			}
			gapLen := size - gapHead - gapTail
			if ringLen := ring.Len(); ringLen != gapLen {
				t.Errorf("Length should be %d, got %d", gapLen, ringLen)
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			if ring.Full() {
				t.Error("Full should be false")
			}
			save := ring
			iter := ring.Iter()
			for i := 0; i < gapLen; i++ {
				if !iter.Next() {
					t.Error("Iterator should not be empty")
				}
				if pos := iter.Pos(); pos != gapHead+i {
					t.Errorf("Iteration number %d should yield %d, got %d", i, gapHead+i, pos)
					t.FailNow()
				}
			}
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}

func TestInvert(t *testing.T) {
	sizes := []int{5, 7, 11} // try a few sizes, just in case
	for _, size := range sizes {
		size := size
		t.Run(fmt.Sprintf("Invert ring of size %d", size), func(t *testing.T) {
			ring := ring.NewRing(size)
			half := size/2 + 1
			for i := 0; i < half; i++ {
				if pos := ring.Push(); pos != i {
					t.Errorf("Index after %d pushes should be %d, got %d", i+1, i, pos)
					t.FailNow()
				}
			}
			for i := 0; i < half-1; i++ {
				if pos := ring.Pop(); pos != i {
					t.Errorf("Index after %d pops should be %d, got %d", i+1, i, pos)
					t.FailNow()
				}
			}
			for i := 0; i < half+1; i++ {
				if pos := ring.Push(); pos != (half+i)%size {
					t.Errorf("Index after %d pops should be %d, got %d", i+1, (half + i), pos)
					t.FailNow()
				}
			}
			gapLen := half + 2
			if ringLen := ring.Len(); ringLen != gapLen {
				t.Errorf("Length should be %d, got %d", gapLen, ringLen)
			}
			if ring.Cap() != size {
				t.Errorf("Size should be %d", size)
			}
			// When size == 5, the process above will fill the buffer
			if gapLen == ring.Cap() {
				if !ring.Full() {
					t.Error("Full should be true")
				}
			} else if ring.Full() {
				t.Error("Full should be false")
			}
			save := ring
			iter := ring.Iter()
			half--
			for i := 0; i < gapLen; i++ {
				if !iter.Next() {
					t.Error("Iterator should not be empty")
				}
				if pos := iter.Pos(); pos != (half+i)%size {
					t.Errorf("Iteration number %d should yield %d, got %d", i, (half+i)%size, pos)
					t.FailNow()
				}
			}
			if iter.Next() {
				t.Error("Iterator should be exhausted")
			}
			if ring != save {
				t.Error("Iterator should not modify ring")
			}
		})
	}
}
