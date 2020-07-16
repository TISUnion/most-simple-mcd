package plugins

import (
	"github.com/TISUnion/most-simple-mcd/modules"
)

func RegisterPlugin() {
	pc := modules.GetPluginContainerInstance()
	//注入插件容器
	for _, v := range plugins {
		pc.RegisterPlugin(v)
	}
}
