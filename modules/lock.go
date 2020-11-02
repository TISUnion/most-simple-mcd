package modules

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"sync"
)

type MLock struct {
	lock    sync.Locker
	isLock  bool
	message string
}

func (l *MLock) IsLock() bool {
	return l.isLock
}

func (l *MLock) Lock() {
	l.lock.Lock()
	l.isLock = true
}

func (l *MLock) Unlock() {
	l.lock.Unlock()
	l.isLock = false
}

func (l *MLock) LockWithMessage(msg string) {
	l.Lock()
	l.message = msg
}

func (l *MLock) GetLockMessage() (msg string) {
	if l.isLock {
		msg = l.message
	}
	return
}

func GetLock() _interface.Lock {
	return &MLock{
		lock: &sync.Mutex{},
	}
}
