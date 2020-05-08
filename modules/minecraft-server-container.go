package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/utils"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var minecraftServerContainer *MinecraftServerContainer

var (
	NO_SERVER  = errors.New("id没有对应的服务器")
	REPEAT_ID  = errors.New("id段冲突")
	NOT_MIRROR = errors.New("不是镜像服务端")
)

type MinecraftServerContainer struct {

	// 所有mc服务器实例
	minecraftServers map[string]server.MinecraftServer

	groupLock *sync.WaitGroup

	// 开启的mc服务器实例
	startServers map[string]server.MinecraftServer

	// 关闭的mc服务器实例
	stopServers map[string]server.MinecraftServer

	// 操作锁
	lock *sync.Mutex

	// 各时间段的回调
	mcCallbacks map[string][]func(string)
}

// 统一注册关闭回调
func (m *MinecraftServerContainer) RegisterAllServerCloseCallback(f func(string)) {
	m.mcCallbacks[constant.MC_CLOSE_CALLBACK] = append(m.mcCallbacks[constant.MC_CLOSE_CALLBACK], f)
	for _, mcS := range m.minecraftServers {
		mcS.RegisterCloseCallback(f)
	}
}

// 统一注册开启回调
func (m *MinecraftServerContainer) RegisterAllServerOpenCallback(f func(string)) {
	m.mcCallbacks[constant.MC_OPEN_CALLBACK] = append(m.mcCallbacks[constant.MC_OPEN_CALLBACK], f)
	for _, mcS := range m.minecraftServers {
		mcS.RegisterOpenCallback(f)
	}
}

// 统一注册保存回调
func (m *MinecraftServerContainer) RegisterAllServerSaveCallback(f func(string)) {
	m.mcCallbacks[constant.MC_SAVE_CALLBACK] = append(m.mcCallbacks[constant.MC_SAVE_CALLBACK], f)
	for _, mcS := range m.minecraftServers {
		mcS.RegisterSaveCallback(f)
	}
}

func (m *MinecraftServerContainer) ChangeConfCallBack() {

}

func (m *MinecraftServerContainer) DestructCallBack() {

}

func (m *MinecraftServerContainer) InitCallBack() {
	m.loadDbServer()
	m.loadLocalServer()
	m.mcCallbacks = make(map[string][]func(string))
}

func (m *MinecraftServerContainer) GetServerById(id string) (server.MinecraftServer, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._getServerById(id)
}

func (m *MinecraftServerContainer) _getServerById(id string) (server.MinecraftServer, error) {
	if len(id) < constant.UUID_LENGTH {
		abId, err := m._getServerLikeId(id)
		if err != nil {
			return nil, err
		}
		id = abId
	}
	if minecraftServer, ok := m.minecraftServers[id]; ok {
		return minecraftServer, nil
	}
	return nil, NO_SERVER
}

// 非完全匹配id
func (m *MinecraftServerContainer) _getServerLikeId(id string) (string, error) {
	aCfg := m._getAllServerConf()
	aRes := make([]string, 1)
	for _, sCfg := range aCfg {
		if strings.Contains(sCfg.EntryId, id) {
			aRes = append(aRes, sCfg.EntryId)
		}
	}
	if len(aRes) > 1 {
		return "", REPEAT_ID
	}
	return aRes[0], nil
}

func (m *MinecraftServerContainer) GetMirrorServerById(id string) (server.MinecraftServer, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	mcS, err := m._getServerById(id)
	if err != nil {
		return nil, err
	}
	if mcS.GetServerConf().IsMirror {
		return mcS, nil
	} else {
		return nil, NOT_MIRROR
	}

}

// 根据id启动服务端
func (m *MinecraftServerContainer) StartById(id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, err := m._getServerById(id)
	if err != nil {
		return err
	}
	err = minecraftServer.Start()
	if err != nil {
		return err
	}
	m.startServers[id] = minecraftServer
	delete(m.stopServers, id)
	return nil
}

func (m *MinecraftServerContainer) StartAll() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	for _, minecraftServer := range m.minecraftServers {
		id := minecraftServer.GetServerConf().EntryId
		// 已经开启则不需要重复启动
		if _, ok := m.startServers[id]; ok {
			continue
		}
		err := minecraftServer.Start()
		if err != nil {
			return err
		}
		m.startServers[id] = minecraftServer
		delete(m.stopServers, id)
	}
	return nil
}

// 根据id停止服务端
func (m *MinecraftServerContainer) StopById(id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, err := m._getServerById(id)
	if err != nil {
		return err
	}
	err = minecraftServer.Stop()
	if err != nil {
		return err
	}
	m.stopServers[id] = minecraftServer
	delete(m.startServers, id)
	return nil
}

// 根据id重启服务端
func (m *MinecraftServerContainer) RestartById(id string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	minecraftServer, err := m._getServerById(id)
	if err != nil {
		return err
	}
	err = minecraftServer.Restart()
	if err != nil {
		return err
	}
	m.startServers[id] = minecraftServer
	delete(m.stopServers, id)
	return nil
}

// 获取所有服务端的配置
func (m *MinecraftServerContainer) GetAllServerConf() []*models.ServerConf {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._getAllServerConf()
}

func (m *MinecraftServerContainer) _getAllServerConf() []*models.ServerConf {
	result := make([]*models.ServerConf, 0)
	for _, v := range m.minecraftServers {
		result = append(result, v.GetServerConf())
	}

	return result
}

// 把根据配置添加服务端
func (m *MinecraftServerContainer) AddServer(config *models.ServerConf, isSave bool) {
	if config.Memory <= 0 {
		config.Memory = constant.MC_DEFAULT_MEMORY
	}
	if len(config.CmdStr) == 0 {
		config.CmdStr = utils.GetCommandArr(config.Memory, config.RunPath)
	}

	mcServer := NewMinecraftServer(config)

	entryId := mcServer.GetServerEntryId()

	// 加入map
	m.minecraftServers[entryId] = mcServer
	m.stopServers[entryId] = mcServer

	// 注册回调
	for _, fCb := range m.mcCallbacks[constant.MC_OPEN_CALLBACK] {
		mcServer.RegisterOpenCallback(fCb)
	}

	for _, fCb := range m.mcCallbacks[constant.MC_CLOSE_CALLBACK] {
		mcServer.RegisterCloseCallback(fCb)
	}

	for _, fCb := range m.mcCallbacks[constant.MC_SAVE_CALLBACK] {
		mcServer.RegisterSaveCallback(fCb)
	}
	if isSave {
		m._saveToDb()
	}
}

// 获取所有服务端对象实例
func (m *MinecraftServerContainer) GetAllServerObj() map[string]server.MinecraftServer {
	return m.minecraftServers
}

// 处理mc服务端文件
func (m *MinecraftServerContainer) HandleMcFile(filePath, name string, port, memory int64) *models.ServerConf {
	path, _ := utils.GetCurrentPath()
	entryId := uuid.NewV4().String()
	_, filename := filepath.Split(filePath)
	// 把服务端jar文件复制到对应文件夹和备份文件夹中，源文件删除
	serverDir := filepath.Join(path, constant.MC_SERVER_DIR, entryId, entryId+".jar")
	serverBackDir := filepath.Join(path, constant.MC_SERVER_BACK, filename)
	copyErr := utils.CopyDir(filePath, serverDir)
	backCopyErr := utils.CopyDir(filePath, serverBackDir)
	deleteErr := os.Remove(filePath)
	if deleteErr != nil {
		WriteLogToDefault(fmt.Sprintf("服务端：%s, 删除失败"), constant.LOG_WARNING)
		return nil
	}
	if copyErr != nil || backCopyErr != nil {
		WriteLogToDefault(fmt.Sprintf("服务端：%s, 复制失败"), constant.LOG_WARNING)
		return nil
	}
	if name != "" {
		filename = name
	}

	// 生成config
	return &models.ServerConf{
		Name:    filename,
		RunPath: serverDir,
		EntryId: entryId,
		Port:    port,
		Memory:  memory,
	}
}

// 读取本地的mc服务端文件
func (m *MinecraftServerContainer) loadLocalServer() {
	m.lock.Lock()
	defer m.lock.Unlock()
	path, _ := utils.GetCurrentPath()
	jarspath, _ := filepath.Glob(fmt.Sprintf("%s/*.jar", path))
	// 读取当前目录下的所有jar文件
	for _, v := range jarspath {
		m.AddServer(m.HandleMcFile(v, "", 0, 0), false)
	}
	m._saveToDb()
}

// 读取数据库中mc配置
func (m *MinecraftServerContainer) loadDbServer() {
	serversConf := m.getServerConfFromDb()
	for _, v := range serversConf {
		// 只有没有被删除的服务端才会加入容器中
		if utils.ExistsResource(v.RunPath) {
			m.AddServer(v, false)
		}
	}
}

// 读取数据库中的服务端配置
func (m *MinecraftServerContainer) getServerConfFromDb() []*models.ServerConf {
	serversConfStr := GetFromDatabase(constant.MC_SERVER_DB_KEY)
	var serversConf []*models.ServerConf
	_ = json.Unmarshal([]byte(serversConfStr), &serversConf)
	// 默认服务端未启动
	for _, c := range serversConf {
		c.State = constant.MC_STATE_STOP
	}
	return serversConf
}

// 持久化服务器配置
func (m *MinecraftServerContainer) SaveToDb() {
	m.lock.Lock()
	defer m.lock.Unlock()
	m._saveToDb()
}

// 持久化服务器配置
func (m *MinecraftServerContainer) _saveToDb() {
	config := m._getAllServerConf()
	data, _ := json.Marshal(config)
	SetFromDatabase(constant.MC_SERVER_DB_KEY, string(data))
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

func GetMinecraftServerContainerInstance() container.MinecraftContainer {
	if minecraftServerContainer != nil {
		return minecraftServerContainer
	}

	minecraftServerContainer = &MinecraftServerContainer{
		minecraftServers: make(map[string]server.MinecraftServer),
		groupLock:        &sync.WaitGroup{},
		startServers:     make(map[string]server.MinecraftServer),
		stopServers:      make(map[string]server.MinecraftServer),
		lock:             &sync.Mutex{},
	}
	RegisterCallBack(minecraftServerContainer)
	return minecraftServerContainer
}
