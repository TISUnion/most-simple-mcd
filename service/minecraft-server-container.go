package service

import (
	"errors"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"sync"
)

var minecraftServerContainer *MinecraftServerContainer

var (
	NO_SERVER = errors.New("id没有对应的服务器")
)

type MinecraftServerContainer struct {

	// 所有mc服务器实例
	minecraftServers map[int]server.MinecraftServer

	groupLock *sync.WaitGroup

	// 开启的mc服务器实例
	startServers map[int]server.MinecraftServer

	// 关闭的mc服务器实例
	stopServers map[int]server.MinecraftServer

	// 操作锁
	lock *sync.Mutex
}

func (m *MinecraftServerContainer) ChangeConfCallBack() {

}

func (m *MinecraftServerContainer) DestructCallBack() {

}

func (m *MinecraftServerContainer) InitCallBack() {

}

func (m *MinecraftServerContainer) GetServerById(id int) (server.MinecraftServer, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._getServerById(id)
}

func (m *MinecraftServerContainer) _getServerById(id int) (server.MinecraftServer, bool) {
	if minecraftServer, ok := m.minecraftServers[id]; ok {
		return minecraftServer, ok
	}
	return nil, false
}

func (m *MinecraftServerContainer) GetMirrorServerById(id int) (server.MinecraftServer, bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._getServerById(id)
}

func (m *MinecraftServerContainer) StartById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	return minecraftServer.Start()
}

func (m *MinecraftServerContainer) StopById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	return minecraftServer.Stop()
}

func (m *MinecraftServerContainer) RestartById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	return minecraftServer.Restart()
}

func (m *MinecraftServerContainer) GetAllServerConf() []*json_struct.ServerConf {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]*json_struct.ServerConf, 0)
	for _, v := range m.minecraftServers{
		result = append(result, v.GetServerConf())
	}

	return result
}

func (m *MinecraftServerContainer) Add(string, server.MinecraftServer) error {
	return nil
}

func (m *MinecraftServerContainer) StopAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, v := range m.minecraftServers{
		_ = v.Stop()
	}
	return nil
}

func GetMinecraftServerContainerObj() container.MinecraftContainer {
	if minecraftServerContainer != nil {
		return minecraftServerContainer
	}

	minecraftServerContainer = &MinecraftServerContainer{
		minecraftServers: make(map[int]server.MinecraftServer),
		groupLock:        &sync.WaitGroup{},
		startServers:     make(map[int]server.MinecraftServer),
		stopServers:      make(map[int]server.MinecraftServer),
		lock:             &sync.Mutex{},
	}

	RegisterCallBack(minecraftServerContainer)

	return minecraftServerContainer
}
