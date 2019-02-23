package queue

import "sync"

// Bytes is a thread-safe queue
type Bytes struct {
	first *nodeBytes
	last  *nodeBytes
	size  uint
	lock  *sync.Mutex
}

// NewBytes returns a new Bytes
func NewBytes() *Bytes {
	return &Bytes{lock: &sync.Mutex{}}
}

// IsEmpty determins if Bytes is empty
func (b *Bytes) IsEmpty() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.first == nil
}

// Size returns the size of Bytes
func (b *Bytes) Size() uint {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.size
}

// Enqueue adds data to Bytes
func (b *Bytes) Enqueue(item []byte) {
	b.lock.Lock()
	defer b.lock.Unlock()

	oldLast := b.last
	b.last = &nodeBytes{item: item}
	if b.IsEmpty() {
		b.first = b.last
	} else {
		oldLast.next = b.last
	}
	b.size++
}

// Dequeue removes data from Bytes
func (b *Bytes) Dequeue() []byte {
	b.lock.Lock()
	defer b.lock.Unlock()

	item := b.first.item
	b.first = b.first.next
	if b.IsEmpty() {
		b.last = nil
	}
	b.size--
	return item
}

// NodeBytes is a linked list of byte array nodes
type nodeBytes struct {
	item []byte
	next *nodeBytes
}
