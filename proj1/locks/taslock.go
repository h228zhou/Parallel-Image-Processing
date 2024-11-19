package locks

import (
	"runtime"
	"sync/atomic"
)

// TASLock definition
type TASLock struct {
	flag int32
}

// Lock acquires the TASLock.
func (lock *TASLock) Lock() {
	for !atomic.CompareAndSwapInt32(&lock.flag, 0, 1) {
		// Busy-waiting: yield the processor to allow other goroutines to run.
		runtime.Gosched()
	}
}

// Unlock releases the TASLock.
func (lock *TASLock) Unlock() {
	atomic.StoreInt32(&lock.flag, 0)
}
