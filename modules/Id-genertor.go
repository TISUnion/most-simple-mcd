package modules

import "sync"

var (
	increateId = 1
	idLock     = &sync.Mutex{}
)

// 全局唯一id生成器
func GetIncreateId() int {
	idLock.Lock()
	defer idLock.Unlock()
	increateId++
	return increateId
}
