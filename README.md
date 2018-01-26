# mutmux
Mutex Multiplexer for synchronized concurrent accesses to a set of shared resources. The original use case was to handle many go routines performing io operations on a set of many files in undeterministic fashion.

## use
```
// get lock before using shared resource
lock := mmux.GetLock(filename)

// exclusive access to resource
useResource(filename)

// release the lock when you are done
lock.Release()
```
