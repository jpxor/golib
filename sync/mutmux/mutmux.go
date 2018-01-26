// Package mutmux provides the ability to synchronize concurrent
// accesses to a set of resources. Named for "Mutex Multiplexer".
//
// The original use case was to handle many go routines performing io
// operations on a set of many files in undeterministic fashion.
package mutmux

import "sync"

// Mutmux is a mutex multiplexer for synchronizing
// access to a set of shared resources.
type Mutmux struct {
	mutex *sync.Mutex
	xmap  map[string]*sync.Mutex
}

// Lock ensures mutual exclusion to the named resource
// and ensures only the owner of the lock can release it.
type Lock struct {
	resid *string
	xmap  map[string]*sync.Mutex
}

// Init will create and initialize the Mutmux.
func Init() Mutmux {
	return Mutmux{}.init()
}

// init will initialize an existing Mutmux.
func (s Mutmux) init() Mutmux {
	s.mutex = &sync.Mutex{}
	s.xmap = make(map[string]*sync.Mutex)
	return s
}

// GetLock attempts to acquire the lock for the named resource.
// This function will block until the lock is acquired.
func (s Mutmux) GetLock(resid string) Lock {
	lock := s.xmap[resid]
	if lock == nil {
		s.mutex.Lock()
		lock = s.xmap[resid]
		if lock == nil {
			lock = &sync.Mutex{}
			s.xmap[resid] = lock
		}
		s.mutex.Unlock()
	}
	lock.Lock()
	return Lock{
		resid: &resid,
		xmap:  s.xmap,
	}
}

// Release will release the lock
func (lock Lock) Release() {
	lock.xmap[*lock.resid].Unlock()
	lock.resid = nil
	lock.xmap = nil
}
