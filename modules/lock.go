package modules

import "sync"

type Lock struct {
	lock   sync.Locker
	isLock bool
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

func GetLock() sync.Locker {
	return &Lock{
		lock:   GetLock(),
	}
}
