# Package Ring

![GitHub](https://img.shields.io/github/license/rjrivero/ring)
[![Maintainability](https://api.codeclimate.com/v1/badges/b6a4dd9540d7815ffec9/maintainability)](https://codeclimate.com/github/rjrivero/ring/maintainability)
[![Go Report](https://goreportcard.com/badge/github.com/rjrivero/ring)](https://goreportcard.com/badge/github.com/rjrivero/git)
[![Build Status](https://travis-ci.org/rjrivero/ring.svg?branch=master)](https://travis-ci.org/rjrivero/ring)
[![Coverage Status](https://coveralls.io/repos/github/rjrivero/ring/badge.svg?branch=master)](https://coveralls.io/github/rjrivero/ring?branch=master)

Package ring provides support for circular data structures, a.k.a. Rings. Rings can be used to implement type safe circular queues or stacks, backed by fixed-size slices.

For type safety, rings are usually coupled with a fixed size slice, e.g.:

```go
import "github.com/rjrivero/ring"

type MyRing struct {
    ring.Ring
    Buf []MyType
}

func New(size int) MyRing {
    return MyRing{
        ring.Ring: ring.New(size),
        Buf: make([]MyType, size),
    }
}
```

You can use your structure as a LIFO Stack with `Push` and `Pop`:

```go
stack := MyRing(16)

stack.Buf[m.Push()] = value1
stack.Buf[m.Push()] = value2
v := stack.Buf[m.Pop()] // v will be value2
```

And iterate it:

```go
for iter := stack.Ring; iter.Some(); {
    v := stack.Buf[iter.Pop()]
    // do something with v
}
```

To use your structure as a Fifo Queue, just replace `Pop` with `PopFront`:

```go
queue := MyRing(16)

queue.Buf[m.Push()] = value1
queue.Buf[m.Push()] = value2
v := queue.Buf[m.PopFront()] // v will be value1

for iter := queue.Ring; iter.Some(); {
    v := queue.Buf[iter.PopFront()]
    // do something with v
}
```

As you see, calling `Push`, `Pop` or `Popfront` returns the index in the buffer for you to perform the actual operation (set or get) in a type-safe manner.

Calling `Full` before `Push`, you can also check if there is room left in the ring, or you are evicting items. This might be useful to perform some cleanup on evicted items, e.g:

```go
full := queue.Full()
tail := queue.Push()
if full {
    // dispose of the current item before overwriting it, e.g.
    queue.Buf[tail].Close()
}
// Now you can store your new value safely
queue.Buf[tail] = value3
```
