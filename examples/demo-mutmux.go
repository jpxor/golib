package main

import (
	"fmt"
	"time"

	"github.com/jpxor/golib/sync/mutmux"
)

// create and initialize
var mmux = mutmux.Init()

func main() {
	filename := "path/to/file"

	// launch multiple threads all trying to
	// access the same file
	go work(filename, 0)
	go work(filename, 1)
	go work(filename, 2)

	// Sleep until all work is done
	time.Sleep(time.Duration(3)*time.Second + time.Microsecond)
}

func work(filename string, id int) {
	fmt.Println(id, ": waiting for lock")

	// get lock before using shared resource
	lock := mmux.GetLock(filename)

	fmt.Println(id, ": lock acquired")
	useResource(filename)
	fmt.Println(id, ": lock released")

	// release the lock when you are done
	lock.Release()
}

func useResource(path string) {
	// Pretend to do work
	time.Sleep(time.Second)
}
