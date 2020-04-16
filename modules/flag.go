package modules

import (
	"flag"
)

var (
	terminalConfs map[string]*string
)

func InitFlag() map[string]*string {
	if terminalConfs != nil {
		return terminalConfs
	}
	terminalConfs = make(map[string]*string)
	for name, confParam := range DefaultConfParam {
		terminalConfs[name] = flag.String(name, "", confParam.Description)
	}
	flag.Parse()
	return terminalConfs
}
