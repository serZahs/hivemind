package main

import (
	"sync"
	"time"
)

type message struct {
	id        int
	payload   []byte
	timestamp time.Time
	priority  int
}

type queue struct {
	lock sync.Mutex
	head int
	tail int
	data []*message
}

const queue_starting_size = 2 // Power of two for efficient modulus

func (q *queue) init() {
	q.head = 0
	q.tail = 0
	q.data = make([]*message, queue_starting_size)
}

func (q *queue) enqueue(msg *message) {
	if (q.tail == q.head) && (q.data[q.tail] != nil) {
		q.grow()
	}
	q.data[q.tail] = msg
	q.tail = q.next(q.tail)
}

func (q *queue) dequeue() *message {
	q.lock.Lock()
	defer q.lock.Unlock()
	if (q.head == q.tail) && (q.data[q.head] == nil) {
		return nil
	}
	msg := q.data[q.head]
	q.data[q.head] = nil
	q.head = q.next(q.head)
	return msg
}

func (q *queue) grow() {
	new_size := len(q.data) * 2 // Power of two!
	new_data := make([]*message, new_size)
	if q.tail > q.head {
		copy(new_data, q.data[q.head:q.tail])
	} else {
		count := copy(new_data, q.data[q.head:])
		copy(new_data[count:], q.data[:q.tail])
	}
	q.head = 0
	q.tail = len(q.data)
	q.data = new_data
}

func (q *queue) next(i int) int {
	return (i + 1) & (len(q.data) - 1)
}
