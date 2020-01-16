package service

import "flag"

func InitFlag() (terminalConfs TerminalType) {
	terminalConfs = make(TerminalType)
	for _, v := range DefaultConfKeys {
		terminalConfs[v] = flag.String(v,"","")
	}
	flag.Parse()
	return
}
