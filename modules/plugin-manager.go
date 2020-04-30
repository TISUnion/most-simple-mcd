package modules

import (
	plugin_interface "github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/models"
	"sync"
)

type PluginManager struct {
	// 所有插件
	allPlugins map[string]plugin_interface.Plugin
	// 可用插件
	ablePlugins map[string]plugin_interface.Plugin
	// 不可用插件
	disablePlugins map[string]plugin_interface.Plugin
	// mc服务端
	mcServ server.MinecraftServer
	// 锁
	lock *sync.Mutex
}

func (m *PluginManager) ChangeConfCallBack() {
}

func (m *PluginManager) DestructCallBack() {
}

func (m *PluginManager) InitCallBack() {
	m.disablePlugins = make(map[string]plugin_interface.Plugin)
	m.lock = &sync.Mutex{}
}

func (m *PluginManager) GetAblePlugins() map[string]plugin_interface.Plugin {
	return m.ablePlugins
}

func (m *PluginManager) GetDisablePlugins() map[string]plugin_interface.Plugin {
	return m.disablePlugins
}

func (m *PluginManager) BanPlugin(pluginId string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.ablePlugins[pluginId]
	if !ok {
		return
	}
	delete(m.ablePlugins, pluginId)
	m.disablePlugins[pluginId] = m.allPlugins[pluginId]
	if !m.allPlugins[pluginId].IsGlobal() {
		m.allPlugins[pluginId].Stop()
	}
}

func (m *PluginManager) UnbanPlugin(pluginId string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.disablePlugins[pluginId]
	if !ok {
		return
	}
	delete(m.disablePlugins, pluginId)
	m.ablePlugins[pluginId] = m.allPlugins[pluginId]
	if !m.allPlugins[pluginId].IsGlobal() {
		m.allPlugins[pluginId].Start()
	}
}

func (m *PluginManager) HandleMessage(msg *json_struct.ReciveMessage) {
	for _, p := range m.ablePlugins {
		p.HandleMessage(msg)
	}
}

// 动态添加插件
func (m *PluginManager) AddPlugin(p plugin_interface.Plugin) {
	m.lock.Lock()
	defer m.lock.Unlock()
	addPname := p.GetName()
	for _, pl := range m.allPlugins {
		if pl.GetName() == addPname {
			return
		}
	}
	p.Init(m.mcServ)
	m.allPlugins[p.GetId()] = p
	m.ablePlugins[p.GetId()] = p
}
