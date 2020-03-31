package modules

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"sync"
)

var pluginContainerObj *PluginContainer

type PluginContainer struct {
	// 全局插件
	globalPlugins map[string]plugin_interface.Plugin

	// 单服务端插件（一个服务端会对应一个实例）
	plugins map[string]plugin_interface.Plugin

	// 锁
	lock *sync.Mutex

	// 插件管理器（动态增加插件）
	managers []plugin_interface.PluginManager
}

func (c *PluginContainer) ChangeConfCallBack() {
}

func (c *PluginContainer) DestructCallBack() {
}

func (c *PluginContainer) InitCallBack() {
	c.globalPlugins = make(map[string]plugin_interface.Plugin)
	c.plugins = make(map[string]plugin_interface.Plugin)
	c.managers = make([]plugin_interface.PluginManager, 0)
	c.lock = &sync.Mutex{}
}

func (c *PluginContainer) RegisterPlugin(p plugin_interface.Plugin) {
	c.lock.Lock()
	defer c.lock.Unlock()
	newPname := p.GetName()
	for _, pl := range c.globalPlugins {
		if pl.GetName() == newPname {
			return
		}
	}
	for _, pl := range c.plugins {
		if pl.GetName() == newPname {
			return
		}
	}
	// 全局插件判断
	if p.IsGlobal() {
		c.globalPlugins[p.GetId()] = p
	} else {
		c.plugins[p.GetId()] = p
	}

	// 分发给各插件管理器
	for _, m := range c.managers {
		var newPl plugin_interface.Plugin
		if !p.IsGlobal() {
			newPl = p.NewInstance()
		}
		m.AddPlugin(newPl)
	}
}

func (c *PluginContainer) NewPluginManager(mcServer server.MinecraftServer) plugin_interface.PluginManager {
	c.lock.Lock()
	defer c.lock.Unlock()
	allplugins := make(map[string]plugin_interface.Plugin)
	for _, p := range c.globalPlugins {
		allplugins[p.GetId()] = p
	}

	for _, p := range c.plugins {
		newP := p.NewInstance()
		allplugins[newP.GetId()] = newP
		newP.Init(mcServer)
	}
	pm := &PluginManager{
		allPlugins:  allplugins,
		ablePlugins: allplugins,
		mcServ:      mcServer,
	}
	RegisterCallBack(pm)
	c.managers = append(c.managers, pm)
	return pm
}

func GetPluginContainerInstance() plugin_interface.PluginContainer {
	if pluginContainerObj != nil {
		return pluginContainerObj
	}
	pc := &PluginContainer{}
	RegisterCallBack(pc)
	pluginContainerObj = pc
	return pc
}
