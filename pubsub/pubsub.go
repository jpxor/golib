// Package pubsub provides a publish-subscribe service for messaging within the same process.
package pubsub

import "container/list"
import "sync"

// Relay provides a the service for all
// pub-sub needs.
type Relay struct {
	rmap     map[string]*list.List
	rwmutex  *sync.RWMutex
	capacity int
}

// Init will create a Relay and set its message
// queue capacity.
func Init(capacity int) Relay {
	return Relay{
		rmap:     make(map[string]*list.List),
		rwmutex:  &sync.RWMutex{},
		capacity: capacity,
	}
}

// Subscribe will provide a channel (chan interface{}) which
// will receive all published messages to the specified key (topic).
// The second return value is the "receiver id" and is used for
// unsubscribing to the key.
func (relay *Relay) Subscribe(key string) (chan interface{}, *list.Element) {
	relay.rwmutex.Lock()
	defer relay.rwmutex.Unlock()
	li := relay.rmap[key]
	if li == nil {
		li = &list.List{}
		relay.rmap[key] = li
	}
	recv := make(chan interface{}, relay.capacity)
	return recv, li.PushBack(recv)
}

// Unsubscribe will remove and close the receiving channel
// identified by the receiver id (rid).
func (relay *Relay) Unsubscribe(key string, rid *list.Element) {
	relay.rwmutex.Lock()
	defer relay.rwmutex.Unlock()
	li := relay.rmap[key]
	if li == nil {
		return
	}
	li.Remove(rid)
	if li.Len() == 0 {
		delete(relay.rmap, key)
	}
	close(rid.Value.(chan interface{}))
}

// Publish will send any data to subscribers via the channel.
// If a subscriber's channel buffer is full, it will be skipped.
func (relay *Relay) Publish(key string, data interface{}) {
	subscribers := relay.copySubList(key)
	for e := subscribers.Front(); e != nil; e = e.Next() {
		ch, ok := e.Value.(chan interface{})
		if ok {
			select {
			case ch <- data: // data sent
			default: // prevents blocking
			}
		}
	}
}

// The list of subscriber channels is copied out for
// better concurrency. This allows the map to be updated
// concurrently with messages being transmited. The copy
// is significantly faster than channel transmit, see
// benchmark: bench-pubsub.go
func (relay *Relay) copySubList(key string) list.List {
	relay.rwmutex.RLock()
	defer relay.rwmutex.RUnlock()
	li := relay.rmap[key]
	if li == nil {
		return list.List{}
	}
	return *li
}
