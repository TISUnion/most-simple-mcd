package service

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
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

// 根据id启动服务端
func (m *MinecraftServerContainer) StartById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	err := minecraftServer.Start()
	if err != nil {
		return err
	}
	m.startServers[id] = minecraftServer
	delete(m.stopServers, id)
	return nil
}

// 根据id停止服务端
func (m *MinecraftServerContainer) StopById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	err := minecraftServer.Stop()
	if err != nil {
		return err
	}
	m.stopServers[id] = minecraftServer
	delete(m.startServers, id)
	return nil
}

// 根据id重启服务端
func (m *MinecraftServerContainer) RestartById(id int) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, ok := m._getServerById(id)
	if !ok {
		return NO_SERVER
	}
	err := minecraftServer.Restart()
	if err != nil {
		return err
	}
	m.startServers[id] = minecraftServer
	delete(m.stopServers, id)
	return nil
}

// 获取所有服务端的配置
func (m *MinecraftServerContainer) GetAllServerConf() []*json_struct.ServerConf {
	m.lock.Lock()
	defer m.lock.Unlock()
	result := make([]*json_struct.ServerConf, 0)
	for _, v := range m.minecraftServers {
		result = append(result, v.GetServerConf())
	}

	return result
}

// 把根据配置添加服务
func (m *MinecraftServerContainer) Add(config *json_struct.ServerConf) {
	if config.Memory <= 0 {
		config.Memory = 1024
	}
	config.CmdStr = utils.GetCommandArr(config.Memory, config.RunPath)

	mcServer := NewMinecraftServer(config)

	entryId := mcServer.GetServerEntryId()
	// 加入map
	m.minecraftServers[entryId] = mcServer
	m.stopServers[entryId] = mcServer

}

// 读取本地的mc服务端文件
func (m *MinecraftServerContainer) loadLocalServer() {
	path, _ := utils.GetCurrentPath()
	jarspath, _ := filepath.Glob(fmt.Sprintf("%s/*.jar", path))
	// 读取当前目录下的所有jar文件
	for _, v := range jarspath {
		_, filename := filepath.Split(v)
		filemd5 := utils.Md5(filename + strconv.FormatInt(time.Now().UnixNano(), 10))
		// 把服务端jar文件复制到对应文件夹和备份文件夹中，源文件删除
		serverDir := filepath.Join(path, constant.MC_SERVER_DIR, filemd5, filemd5+".jar")
		serverBackDir := filepath.Join(path, constant.MC_SERVER_BACK, filename)
		copyErr := utils.CopyDir(v, serverDir)
		backCopyErr := utils.CopyDir(v, serverBackDir)
		deleteErr := os.Remove(v)
		if deleteErr != nil {
			GetLogContainerInstance().WriteLog(fmt.Sprintf("服务端：%s, 删除失败"), constant.LOG_WARNING)
			continue
		}
		if copyErr != nil || backCopyErr != nil {
			GetLogContainerInstance().WriteLog(fmt.Sprintf("服务端：%s, 复制失败"), constant.LOG_WARNING)
			continue
		}
		// 生成config
		config := &json_struct.ServerConf{
			Name:    filename,
			RunPath: serverDir,
		}
		m.Add(config)
	}
}

// 停止所有服务端
func (m *MinecraftServerContainer) StopAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, v := range m.minecraftServers {
		_ = v.Stop()
		id := v.GetServerEntryId()
		m.stopServers[id] = v
		delete(m.startServers, id)
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
