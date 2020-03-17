package modules

import (
	"flag"
	"sync"
)

var onceLock = &sync.Once{}

func InitFlag() (terminalConfs map[string]*string) {
	onceLock.Do(func() {
		terminalConfs = make(map[string]*string)
		for name, confParam := range DefaultConfParam {
			terminalConfs[name] = flag.String(name, "", confParam.Description)
		}
		flag.Parse()
	})
	return
}
