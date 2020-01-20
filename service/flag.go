package service

import (
	"flag"
	"sync"
)

var onceLock *sync.Once

func init()  {
	onceLock = &sync.Once{}
}

func InitFlag() (terminalConfs TerminalType) {
	onceLock.Do(func() {
		terminalConfs = make(TerminalType)
		for _, v := range DefaultConfKeys {
			terminalConfs[v] = flag.String(v,"","")
		}
		flag.Parse()
	})
	return
}
