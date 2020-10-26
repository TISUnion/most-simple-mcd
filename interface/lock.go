package _interface

import "sync"

type Lock interface {
	sync.Locker
	LockWithMessage(string)
	GetLockMessage() string // 只有已加锁时才会返回加锁信息
	IsLock() bool
}
