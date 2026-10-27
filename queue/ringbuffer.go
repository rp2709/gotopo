package queue

import "iter"

type RingBuffer[T any] struct {
	buf        []T
	head, tail int
	size       int
}

func NewRingBuffer[T any](initialCapacity int) *RingBuffer[T] {
	if initialCapacity < 0 {
		panic("ringbuffer: capacity cannot be negative")
	}
	return &RingBuffer[T]{
		buf: make([]T, initialCapacity),
	}
}

// Push appends v, growing the buffer if necessary.
func (r *RingBuffer[T]) Push(v T) {
	if r.size == len(r.buf) {
		r.grow()
	}

	r.buf[r.tail] = v
	r.tail++
	if r.tail == len(r.buf) {
		r.tail = 0
	}
	r.size++
}

// grow doubles capacity and preserves FIFO order.
func (r *RingBuffer[T]) grow() {
	old := r.buf

	newCapacity := len(old) * 2
	if newCapacity == 0 {
		newCapacity = 1
	}

	r.buf = make([]T, newCapacity)

	// Copy elements in logical order, handling wraparound.
	n := copy(r.buf, old[r.head:])
	copy(r.buf[n:], old[:r.head])

	r.head = 0
	r.tail = r.size
}

// Pop removes and returns the oldest element.
func (r *RingBuffer[T]) Pop() T {
	if r.size == 0 {
		panic("pop on empty ringbuffer")
	}

	v := r.buf[r.head]

	// Release references to the removed value.
	var zero T
	r.buf[r.head] = zero

	r.head++
	if r.head == len(r.buf) {
		r.head = 0
	}
	r.size--

	return v
}

// Peek returns the oldest element without removing it.
func (r *RingBuffer[T]) Peek() (T, bool) {
	if r.size == 0 {
		var zero T
		return zero, false
	}
	return r.buf[r.head], true
}

func (r *RingBuffer[T]) Len() int {
	return r.size
}

func (r *RingBuffer[T]) Capacity() int {
	return len(r.buf)
}

func (r *RingBuffer[T]) Empty() bool {
	return r.size == 0
}

// Clear empties the buffer and releases references.
func (r *RingBuffer[T]) Clear() {
	var zero T
	for i := range r.buf {
		r.buf[i] = zero
	}
	r.head = 0
	r.tail = 0
	r.size = 0
}

func (r *RingBuffer[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		buf := r.buf
		size := r.size
		i := r.head

		for range size {
			if !yield(buf[i]) {
				return
			}

			i++
			if i == len(buf) {
				i = 0
			}
		}
	}
}
