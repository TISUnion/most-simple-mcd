package service

import (
	"flag"
	"sync"
)

var onceLock *sync.Once

func init() {
	onceLock = &sync.Once{}
}

func InitFlag() (terminalConfs TerminalType) {
	onceLock.Do(func() {
		terminalConfs = make(TerminalType)
		for name, confParam := range DefaultConfParam {
			terminalConfs[name] = flag.String(name, "", confParam.Description)
		}
		flag.Parse()
	})
	return
}
