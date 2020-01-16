package utils

import (
	"gopkg.in/ini.v1"
)

func SetIniCfg(data map[string]string) *ini.File {
	cfg := ini.Empty()
	sec, _ := cfg.NewSection("")
	for k, v := range data {
		_, _ = sec.NewKey(k, v)
	}
	return cfg
}

func PanicError(msg string) {
	// todo 写入日志

	panic(msg)
}