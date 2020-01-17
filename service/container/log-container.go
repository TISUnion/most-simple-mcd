package container

import (
	"github.com/TISUnion/most-simple-mcd/service"
	"sync"
)

var (
	increateId int         = 0
	idLock     *sync.Mutex = &sync.Mutex{}
)

func getIncreateId() int {
	idLock.Lock()
	defer idLock.Unlock()
	increateId++
	return increateId
}

type LogContainer struct {
	NameIdMapping map[string]int
	Logs          map[int]*service.Log
}
