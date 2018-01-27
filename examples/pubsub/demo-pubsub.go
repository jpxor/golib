package main

import "fmt"
import "time"
import "github.com/jpxor/golib/pubsub"

// Message is an arbitrary struct
type Message struct {
	ID    string
	Value interface{}
}

func main() {

	// Create the pubsub instance
	relay := pubsub.Init(4)

	// Subscribe to "test_channel"
	recv, rid := relay.Subscribe("test_channel")

	// Unsubscribe at some point in the future
	go func() {
		time.Sleep(2 * time.Second)
		relay.Unsubscribe("test_channel", rid)
	}()

	// Publish a message immediately
	relay.Publish("test_channel", &Message{
		ID:    "send anything!",
		Value: []byte("such as a pointer for transfering control of a large dataset without copying any of it"),
	})

	// Publish a message in the future
	go func() {
		time.Sleep(time.Second)
		relay.Publish("test_channel", &Message{
			ID:    "second:",
			Value: []byte("this was sent 1 second later"),
		})
	}()

	// Recieve messages from channel
	for {
		val := <-recv
		if val == nil {
			fmt.Printf("reciever was closed\n")
			break
		}
		msgPtr := val.(*Message)
		fmt.Printf("%s\n", msgPtr)
	}

}
