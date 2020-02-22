# Package Ring

[![Maintainability](https://api.codeclimate.com/v1/badges/b6a4dd9540d7815ffec9/maintainability)](https://codeclimate.com/github/rjrivero/ring/maintainability)
[![Go Report](https://goreportcard.com/badge/github.com/rjrivero/ring)](https://goreportcard.com/badge/github.com/rjrivero/git)
[![Build Status](https://travis-ci.org/rjrivero/ring.svg?branch=master)](https://travis-ci.org/rjrivero/ring)
[![Coverage Status](https://coveralls.io/repos/github/rjrivero/ring/badge.svg?branch=master)](https://coveralls.io/github/rjrivero/ring?branch=master)

Package ring provides support for circular queues backed by fixed-size buffers, in a type-safe manner. E.g.:

```go
import "github.com/rjrivero/ring"

// NodeQueue is a type-safe circular FIFO queue of Node objects
type NodeQueue struct {
  Ring
  Buffer []Node
}

// Create a queue of a given size
func NewQueue(size int) NodeQueue {
    return NodeQueue{
        Ring: ring.NewRing(size),
        Buffer: make([]Node, size)
    }
}

myQueue := NewQueue(16)

// Push a Node to the queue
myQueue.Buffer[myQueue.Push()] = myNode

// Pop a Node from the Queue.
if myQueue.Len() > 0 {
    popped := myQueue.Buffer[myQueue.Pop()]
}

// Iterate over the items in the queue
for iter := myQueue.Iter(); iter.Next(); {
    current := myQueue.Buffer[iter.Pos()]
}
```

Ring is not concurrency-safe, it should not be manipulated by different goroutines without proper synchronization.
