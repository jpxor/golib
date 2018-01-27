// Package main for benchmarking pubsub:
// In PubSub, during a Publish, the list of subscriber
// channels is copied out to reduce readlock time and
// allow for better concurrency. This lets the map be
// updated concurrently with messages being transmited.
//
// This benchmark shows that the copy is significantly
// faster than channel transmit.
package main

import "fmt"
import "time"
import "container/list"

func main() {

	count := 19999999
	li := list.List{}

	// Create the list of channels
	start := time.Now()
	for i := 0; i < count; i++ {
		li.PushBack(make(chan interface{}, 10))
	}
	elapsed := time.Since(start)
	fmt.Printf("list creation took %s\n", elapsed)

	// copy the list to new list instance
	start = time.Now()
	newlist := li
	elapsed = time.Since(start)
	fmt.Printf("copy took %s\n", elapsed)

	// show that the copy succeeded and created a
	// separate instance.
	newlist.PushBack(make(chan interface{}, 10))
	fmt.Printf("newlist size: %+v\n", newlist.Len())
	fmt.Printf("original list size: %+v\n", li.Len())

	// send a 64 bit value into the channel
	start = time.Now()
	for e := li.Front(); e != nil; e = e.Next() {
		ch, ok := e.Value.(chan interface{})
		if ok {
			select {
			case ch <- uint64(101010101): // data sent
			default: // prevents blocking
			}
		}
	}
	elapsed = time.Since(start)
	fmt.Printf("chan transmit took %s\n", elapsed)

	// try again with simpler mechanics
	start = time.Now()
	for e := li.Front(); e != nil; e = e.Next() {
		ch, ok := e.Value.(chan interface{})
		if ok {
			ch <- uint8(0)
		}
	}
	elapsed = time.Since(start)
	fmt.Printf("second chan transmit took %s\n", elapsed)
}
