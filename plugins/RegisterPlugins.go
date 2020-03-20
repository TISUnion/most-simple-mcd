package plugins

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/plugins/reread-chicken"
)

var plugins = []plugin_interface.Plugin{
	reread_chicken.RereadChickenPluginObj,
}

func RegisterPlugin() {
	pc := modules.GetPluginContainerInstance()
	//注入插件容器
	for _, v := range plugins {
		pc.RegisterPlugin(v)
	}
}
