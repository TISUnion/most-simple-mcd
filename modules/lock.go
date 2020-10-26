package modules

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"sync"
)

type Lock struct {
	lock    sync.Locker
	isLock  bool
	message string
}

func (l *Lock) IsLock() bool {
	return l.isLock
}

func (l *Lock) Lock() {
	l.lock.Lock()
	l.isLock = true
}

func (l *Lock) Unlock() {
	l.lock.Unlock()
	l.isLock = false
}

func (l *Lock) LockWithMessage(msg string) {
	l.Lock()
	l.message = msg
}

func (l *Lock) GetLockMessage() (msg string) {
	if l.isLock {
		msg = l.message
	}
	return
}

func GetLock() _interface.Lock {
	return &Lock{
		lock: GetLock(),
	}
}
