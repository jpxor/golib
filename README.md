# jpxor Golib
Library of useful Go packages. 
- mutmux
- pubsub

## Mutmux
Mutex Multiplexer for synchronized concurrent accesses to a set of shared resources. The original use case was to handle many go routines performing io operations on a set of many files in undeterministic fashion.

### use
```
// create Mutmux
var mmux = mutmux.Init()

// get lock before using shared resource
lock := mmux.GetLock(filename)

// exclusive access to resource
useResource(filename)

// release the lock when you are done
lock.Release()
```

## Pubsub
Publish-Subscribe messaging service.

### use
```
// Create the pubsub instance
ps := pubsub.Init(4)

// Subscribe to named channel (topic), 
// Provides a receiving chan and an id
// for unsubscribing later.
recv, rid := ps.Subscribe("demo")

// Publish any message
ps.Publish("demo", "Hello PubSub!")
ps.Publish("demo", &ArbitraryStruct)

// Recieve messages like you would
// with any Go chan
for {
   val := <-recv
   if val == nil {
      fmt.Printf("reciever was closed\n")
      break
   }
   msgPtr := val.(*Message)
   fmt.Printf("%s\n", msgPtr)
}

// Unsubscribe using receiver id
relay.Unsubscribe("demo", rid)
```
