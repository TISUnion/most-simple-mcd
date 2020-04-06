package modules

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io"
	"net"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
)

var (
	PORT_REPEAT_ERROR = errors.New("服务器端口已被其他程序占用，请更换端口或者开启自动更换端口")
)

// MinecraftServer
// mc服务器子进程对象
type MinecraftServer struct {
	*json_struct.ServerConf

	// CmdObj
	//子进程实例
	CmdObj *exec.Cmd

	// stdin
	// 用于关闭输入管道
	stdin io.WriteCloser

	// stdout
	// 子进程输出
	stdout io.ReadCloser

	// Pid
	// 进程pid
	Pid int

	// lock
	// 输入管道同步锁
	lock *sync.Mutex

	// messageChan
	// 玩家发言存储chan
	messageChan chan *json_struct.ReciveMessage

	// MonitorServer
	// 资源监听器
	monitorServer server.MonitorServer

	// 插件管理器
	pluginManager plugin.PluginManager

	// 其他模块订阅服务的消息推送管道
	subscribeMessageChans []chan *json_struct.ReciveMessage

	// 接收关闭服务器信号管道
	stopTagChan chan struct{}

	// Logger
	// 服务端对应日志
	logger _interface.Log
}

func (m *MinecraftServer) BanPlugin(pluginId string) {
	m.pluginManager.BanPlugin(pluginId)
}

func (m *MinecraftServer) UnbanPlugin(pluginId string) {
	m.pluginManager.UnbanPlugin(pluginId)
}

func (m *MinecraftServer) GetPluginsInfo() []*json_struct.PluginInfo {
	res := make([]*json_struct.PluginInfo, 0)
	ablePlugins := m.pluginManager.GetAblePlugins()
	for _, p := range ablePlugins {
		res = append(res, &json_struct.PluginInfo{
			Name:        p.GetName(),
			Id:          p.GetId(),
			IsBan:       false,
			CommandName: p.GetCommandName(),
			Description: p.GetDescription(),
		})
	}
	disablePlugins := m.pluginManager.GetDisablePlugins()
	for _, p := range disablePlugins {
		res = append(res, &json_struct.PluginInfo{
			Name:        p.GetName(),
			Id:          p.GetId(),
			IsBan:       true,
			CommandName: p.GetCommandName(),
			Description: p.GetDescription(),
		})
	}
	return res
}

func (m *MinecraftServer) RegisterSubscribeMessageChan(c chan *json_struct.ReciveMessage) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.subscribeMessageChans = append(m.subscribeMessageChans, c)
}

func (m *MinecraftServer) GetServerEntryId() string {
	return m.EntryId
}

func (m *MinecraftServer) StartMonitorServer() {
	if !m.IsStartMonitor {
		m.IsStartMonitor = true
		m.monitorServer = NewMonitorServer(m.EntryId, m.Pid)
		_ = m.monitorServer.Start()
	}
}

func (m *MinecraftServer) StopMonitorServer() {
	if m.IsStartMonitor {
		m.IsStartMonitor = false
		_ = m.monitorServer.Stop()
	}
}

func (m *MinecraftServer) GetServerMonitor() server.MonitorServer {
	return m.monitorServer
}

func (m *MinecraftServer) GetServerConf() *json_struct.ServerConf {
	return m.ServerConf
}

func (m *MinecraftServer) SetServerConf(c *json_struct.ServerConf) {
	m.ServerConf = c
}

func (m *MinecraftServer) SetMemory(memory int) {
	if memory > 0 {
		m.Memory = memory
	}
}

func (m *MinecraftServer) Rename(name string) {
	if name != "" {
		m.Name = name
	}
}

func (m *MinecraftServer) ChangeConfCallBack() {
}

func (m *MinecraftServer) InitCallBack() {
	// 开启处理接收消息的协成
	go m.handleMessage()

	// 初始化ip
	m.initLocalIps()

	// 创建插件管理器
	m.pluginManager = GetPluginContainerInstance().NewPluginManager(m)

	m.stopTagChan = make(chan struct{})
}

func (m *MinecraftServer) DestructCallBack() {
	_ = m.Stop()
}

// 启动进程
func (m *MinecraftServer) runProcess() error {
	// 校验eula
	if err := m.validateEula(); err != nil {
		return err
	}
	// 校验端口
	if port, err := m.validatePort(); err != nil {
		return err
	} else {
		m.Port = port
	}
	if err := m.CmdObj.Start(); err != nil {
		return err
	}

	// 开启消息
	go m.reciveMessageToChan()

	m.Pid = m.CmdObj.Process.Pid

	if m.IsStartMonitor {
		m.monitorServer = NewMonitorServer(m.EntryId, m.Pid)
		_ = m.monitorServer.Start()
	}
	return nil
}

func (m *MinecraftServer) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	// 重置cmd对象
	m.resetParams()
	// 开启服务端回调
	m.pluginManager.OpenMcServerCallBack()
	if m.State != constant.MC_STATE_STOP {
		return nil
	}
	m.State = constant.MC_STATE_STARTIND
	if err := m.runProcess(); err != nil {
		m.State = constant.MC_STATE_STOP
		return err
	}
	return nil
}

func (m *MinecraftServer) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.State != constant.MC_STATE_START {
		return nil
	}
	if err := m._command("stop"); err != nil {
		WriteLogToDefault(fmt.Sprintf("服务器: %s 关闭失败, 原因：%v", m.Name, err), constant.LOG_ERROR)
		_ = m.CmdObj.Process.Kill()
	}
	m.State = constant.MC_STATE_STOPING
	<-m.stopTagChan
	m.State = constant.MC_STATE_STOP
	// 关闭服务端回调
	m.pluginManager.CloseMcServerCallBack()
	return nil
}

func (m *MinecraftServer) Restart() error {
	if m.State == constant.MC_STATE_START {
		if err := m.Stop(); err != nil {
			return err
		}
	}
	if err := m.Start(); err != nil {
		return err
	}
	return nil
}

// 获取一条消息
func (m *MinecraftServer) resiveOneMessage() ([]byte, error) {
	const MAX_SIZE = 1024
	buff := make([]byte, MAX_SIZE)
	n, err := m.stdout.Read(buff)
	buff = buff[:n]
	if err != nil {
		errMsg := fmt.Sprintf("服务器: %s，已关闭。因为%v", m.Name, err)
		return []byte{}, errors.New(errMsg)
	}
	// 如果一次的数据为1024，就多次获取
	if n == MAX_SIZE {
		for {
			subBuff := make([]byte, MAX_SIZE)
			subN, subErr := m.stdout.Read(buff)
			if subErr != nil {
				errMsg := fmt.Sprintf("服务器: %s，已关闭。因为%v", m.Name, err)
				return []byte{}, errors.New(errMsg)
			}
			subBuff = subBuff[:subN]
			buff = append(buff, subBuff...)
			if subN != MAX_SIZE {
				break
			}
		}
	}
	buff, _ = utils.ParseCharacter(buff)
	return buff, nil
}

// 获取消息，并写入到管道中
func (m *MinecraftServer) reciveMessageToChan() {
	for {
		everyBuff, err := m.resiveOneMessage()
		if err != nil {
			m.State = constant.MC_STATE_STOP
			WriteLogToDefault(err.Error(), constant.LOG_ERROR)
			return
		}
		msg := utils.ParseMessage(everyBuff)
		msg.ServerId = m.EntryId
		m.messageChan <- msg
	}
}

// 处理服务端进程获取的消息
func (m *MinecraftServer) handleMessage() {
	for {
		msg := <-m.messageChan
		//fmt.Print(string(msg.OriginData))
		m.WriteLog(string(msg.OriginData), constant.LOG_INFO)
		if m.Version == "" {
			m.getVersion(msg.OriginData)
		}
		if m.GameType == "" {
			m.getGameType(msg.OriginData)
		}

		// 正在启动
		if m.State == constant.MC_STATE_STARTIND {
			m.sureServerStart(msg.OriginData)
			continue // 如果还没启动，就不分发消息
		}

		// 正在关闭
		if m.State == constant.MC_STATE_STOPING {
			m.sureServerStop(msg.OriginData)
			continue // 如果还没关闭，就不分发消息
		}

		// 分发给插件
		go func() {
			m.pluginManager.HandleMessage(msg)
		}()

		// 分发给各已订阅模块
		go func() {
			for _, c := range m.subscribeMessageChans {
				c <- msg
			}
		}()
	}
}

func (m *MinecraftServer) getVersion(data []byte) {
	reg, _ := regexp.Compile("(([0-9]*\\.[0-9]*\\.{0,1}[0-9]*)+)")
	ves := reg.Find(data)
	if len(ves) > 0 {
		m.Version = string(ves)
	}
}

func (m *MinecraftServer) getGameType(data []byte) {
	reg, _ := regexp.Compile("Default game type: (?P<type>[a-zA-Z]+)")
	match := reg.FindSubmatch(data)
	if len(match) > 1 {
		m.GameType = string(match[1])
	}
}

// 判断服务端是否已经启动
func (m *MinecraftServer) sureServerStart(data []byte) {
	reg, _ := regexp.Compile("\\[Server thread/INFO\\]: Done \\(.*\\)! For help, type \"help\" or \"\\?\"")
	match := reg.Find(data)
	if len(match) > 0 {
		m.State = constant.MC_STATE_START
	}
}

// 判断服务端是否已经关闭
func (m *MinecraftServer) sureServerStop(data []byte) {
	reg, _ := regexp.Compile("\\[Server Shutdown Thread/INFO\\]")
	match := reg.Find(data)
	// 如果已关闭则发送关闭信息
	if len(match) > 0 {
		m.State = constant.MC_STATE_START
		m.stopTagChan <- struct{}{}
	}
}

func (m *MinecraftServer) Command(c string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._command(c)
}

func (m *MinecraftServer) _command(c string) error {
	// 不加换行无法执行命令
	c += "\n"
	_, err := m.stdin.Write([]byte(c))
	return err
}

// validatePort
// 校验mc的端口
func (m *MinecraftServer) validatePort() (int, error) {
	runDir := filepath.Dir(m.RunPath)
	// 新建mc配置文件
	mcConfPath := filepath.Join(runDir, constant.MC_CONF_NAME)
	if f, e := utils.CreateFile(mcConfPath); e == nil {
		f.Close()
	}
	cfg, _ := ini.Load(mcConfPath)
	var realPort int
	if m.Port != 0 {
		realPort = m.Port
		realPortStr := strconv.Itoa(realPort)
		sec, _ := cfg.GetSection(ini.DefaultSection)
		if sec.HasKey(constant.MC_PORT_TEXT) {
			sec.Key(constant.MC_PORT_TEXT).SetValue(realPortStr)
		} else {
			_, _ = sec.NewKey(constant.MC_PORT_TEXT, realPortStr)
		}
		err := cfg.SaveTo(mcConfPath)
		if err != nil {
			WriteLogToDefault(err.Error(), constant.LOG_ERROR)
		}
	} else {
		realPort = constant.MC_DEFAULT_PORT
	}
	// 开启的服务端的端口已被占用,修修改
	if p, _ := utils.GetFreePort(realPort); p == 0 {
		p, err := m.changePort(cfg, mcConfPath, 0)
		if err != nil {
			return 0, err
		}
		realPort = p
	}
	return realPort, nil
}

// changePort
// 更换mc服务端端口
func (m *MinecraftServer) changePort(cfg *ini.File, path string, port int) (int, error) {
	// 如果可以自动更换端口就自动更换端口
	if isChange, _ := strconv.ParseBool(GetConfVal(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT)); isChange {
		unusedPort, _ := utils.GetFreePort(port)
		sec, _ := cfg.GetSection(ini.DefaultSection)
		unusedPortStr := strconv.Itoa(unusedPort)
		// 重新配置文件
		if sec.HasKey(constant.MC_PORT_TEXT) {
			sec.Key(constant.MC_PORT_TEXT).SetValue(unusedPortStr)
		} else {
			_, _ = sec.NewKey(constant.MC_PORT_TEXT, unusedPortStr)
		}
		if err := cfg.SaveTo(path); err != nil {
			return 0, err
		}
		return unusedPort, nil
	} else {
		msg := fmt.Sprintf("服务端：%s，对应的服务器端口已被其他程序占用，请更换端口或者开启自动更换端口", m.Name)
		WriteLogToDefault(msg)
		return 0, PORT_REPEAT_ERROR
	}
}

// validateEula
// 校验mc的eula文件
func (m *MinecraftServer) validateEula() error {
	runDir := filepath.Dir(m.RunPath)
	path := filepath.Join(runDir, constant.EULA_FILE_NAME)
	f, _ := utils.CreateFile(path)
	_ = f.Close()
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}
	sec, err := cfg.GetSection(ini.DefaultSection)
	if err != nil {
		return err
	}
	if sec.HasKey(constant.EULA) {
		eula, err := sec.Key(constant.EULA).Bool()
		if err != nil {
			eula = false
		}
		if !eula {
			sec.Key(constant.EULA).SetValue(constant.TRUE_STR)
			_ = cfg.SaveTo(path)
		}
	} else {
		_, _ = sec.NewKey(constant.EULA, constant.TRUE_STR)
		_ = cfg.SaveTo(path)
	}
	return nil
}

func (m *MinecraftServer) resetParams() {
	if m.stdin != nil {
		_ = m.stdin.Close()
	}
	if m.stdout != nil {
		_ = m.stdout.Close()
	}
	// 重新拼接运行命令
	if len(m.CmdStr) == 0 {
		m.CmdStr = utils.GetCommandArr(m.Memory, m.RunPath)
	}
	m.CmdObj = exec.Command(m.CmdStr[0], m.CmdStr[1:]...)
	m.stdin, _ = m.CmdObj.StdinPipe()
	m.stdout, _ = m.CmdObj.StdoutPipe()
	m.CmdObj.Dir = filepath.Dir(m.RunPath)
	if m.monitorServer != nil {
		// 关闭这个监控器
		m.monitorServer.DestructCallBack()
		m.monitorServer = nil
	}
}

func (m *MinecraftServer) WriteLog(msg string, level string) {
	// 写入自己的日志
	m.logger.WriteLog(&_interface.LogMsgType{
		Message: msg,
		Level:   level,
	})
}

func (m *MinecraftServer) initLocalIps() {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	ips := make([]string, 0)
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	m.Ips = ips
}

// NewMinecraftServer
// 新建一个mc服务端进程
func NewMinecraftServer(serverConf *json_struct.ServerConf) server.MinecraftServer {
	minecraftServer := &MinecraftServer{
		ServerConf:  serverConf,
		lock:        &sync.Mutex{},
		messageChan: make(chan *json_struct.ReciveMessage, 10),
		logger:      GetLogContainerInstance().AddLog(serverConf.EntryId),
	}
	RegisterCallBack(minecraftServer)
	return minecraftServer
}
