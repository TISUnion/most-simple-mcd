package modules

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/plugin"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io"
	"net"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var (
	PORT_REPEAT_ERROR = errors.New("服务器端口已被其他程序占用，请更换端口或者开启自动更换端口")
)

// MinecraftServer
// mc服务器子进程对象
type MinecraftServer struct {
	*models.ServerConf

	*ServerAdapter

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
	messageChan chan string

	// MonitorServer
	// 资源监听器
	monitorServer server.MonitorServer

	// 插件管理器
	pluginManager plugin.PluginManager

	// 接收关闭服务器信号管道
	stopTagChan chan struct{}

	// Logger
	// 服务端对应日志
	logger _interface.Log

	// 其他模块订阅服务的消息推送管道
	subscribeMessageChans []chan *models.ReciveMessage

	// 服务端关闭回调
	mcCloseCallback []func(string)

	// 服务端开启回调
	mcOpenCallback []func(string)

	// 服务端保存回调（执行save-all后调用)
	mcSaveCallback []func(string)
}

func (m *MinecraftServer) BanPlugin(pluginId string) {
	m.pluginManager.BanPlugin(pluginId)
}

func (m *MinecraftServer) UnbanPlugin(pluginId string) {
	m.pluginManager.UnbanPlugin(pluginId)
}

// 获取插件信息
func (m *MinecraftServer) GetPluginsInfo() []*models.PluginInfo {
	res := make([]*models.PluginInfo, 0)
	ablePlugins := m.pluginManager.GetAblePlugins()
	for _, p := range ablePlugins {
		res = append(res, &models.PluginInfo{
			Name:            p.GetName(),
			Id:              p.GetId(),
			IsBan:           false,
			CommandName:     p.GetCommandName(),
			Description:     p.GetDescription(),
			HelpDescription: p.GetHelpDescription(),
		})
	}
	disablePlugins := m.pluginManager.GetDisablePlugins()
	for _, p := range disablePlugins {
		res = append(res, &models.PluginInfo{
			Name:            p.GetName(),
			Id:              p.GetId(),
			IsBan:           true,
			CommandName:     p.GetCommandName(),
			Description:     p.GetDescription(),
			HelpDescription: p.GetHelpDescription(),
		})
	}
	return res
}

func (m *MinecraftServer) RegisterSubscribeMessageChan(c chan *models.ReciveMessage) {
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

func (m *MinecraftServer) GetServerConf() *models.ServerConf {
	return m.ServerConf
}

func (m *MinecraftServer) SetServerConf(c *models.ServerConf) {
	m.ServerConf = c
}

func (m *MinecraftServer) SetMemory(memory int64) {
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

func (m *MinecraftServer) RegisterCloseCallback(c func(string)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.mcCloseCallback = append(m.mcCloseCallback, c)
}

func (m *MinecraftServer) RegisterOpenCallback(c func(string)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.mcOpenCallback = append(m.mcOpenCallback, c)
}

func (m *MinecraftServer) RegisterSaveCallback(c func(string)) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.mcSaveCallback = append(m.mcSaveCallback, c)
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
	// 运行关闭回调
	for _, f := range m.mcCloseCallback {
		f(m.EntryId)
	}
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
func (m *MinecraftServer) resiveOneMessage() (string, error) {
	result := make([]byte, 0)
	for {
		buff := make([]byte, constant.MAX_RESIVE_BUFF_SIZE)
		n, err := m.stdout.Read(buff)
		if err != nil {
			errMsg := fmt.Sprintf("服务器: %s，已关闭。因为%v", m.Name, err)
			m.sureServerStop()
			return "", errors.New(errMsg)
		}
		if n <= 0 {
			continue
		}
		result = append(result, buff[:n]...)
		if buff[n-1] == '\n' {
			break
		}
	}

	result, _ = utils.ParseCharacter(result)
	return string(result), nil
}

// 解析信息
func (m *MinecraftServer) parseMessage(originMsg string) *models.ReciveMessage {
	msgObj, ok := m.GetMessageRegularExpression(originMsg)
	if ok {
		return msgObj
	}
	return &models.ReciveMessage{OriginData: originMsg}
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
		m.messageChan <- everyBuff
	}
}

// 处理服务端进程获取的消息
func (m *MinecraftServer) handleMessage() {
	for {
		msg := <-m.messageChan
		m.WriteLog(msg, constant.LOG_INFO)
		if m.Version == "" {
			m.getVersion(msg)
		}
		if m.GameType == "" {
			m.getGameType(msg)
		}
		msgObj :=  m.parseMessage(msg)
		// 分发给各已订阅模块
		go func() {
			for _, c := range m.subscribeMessageChans {
				c <- msgObj
			}
		}()

		// 分发给插件
		go func() {
			m.pluginManager.HandleMessage(msgObj)
		}()

		// 正在启动
		if m.State == constant.MC_STATE_STARTIND {
			m.sureServerStart(msg)
			continue // 如果还没启动，就不分发消息
		}

		// 正在关闭
		if m.State == constant.MC_STATE_STOPING {
			continue // 如果还没关闭，就不分发消息
		}

		// 如果是保存服务端就调用回调
		go m.sureServerSave(msg)
	}
}

func (m *MinecraftServer) getVersion(data string) {
	if res, ok := m.GetVersionRegularExpression(data); ok {
		m.Version = res.Content
	}
}

func (m *MinecraftServer) getGameType(data string) {
	if res, ok := m.GetVersionRegularExpression(data); ok {
		m.GameType = res.Content
	}
}

// 判断服务端是否已经启动
func (m *MinecraftServer) sureServerStart(data string) {
	_, ok := m.GetGameStartRegularExpression(data)
	if ok {
		m.State = constant.MC_STATE_START
		// 运行开启回调
		for _, f := range m.mcOpenCallback {
			f(m.EntryId)
		}
	}
}

// 判断服务端是否已经关闭
func (m *MinecraftServer) sureServerStop() {
	m.stopTagChan <- struct{}{}
}

// 判断服务端是否已保存
func (m *MinecraftServer) sureServerSave(data string) {
	_, ok := m.GetGameSaveRegularExpression(data)
	// 如果已关闭则发送关闭信息
	if ok {
		for _, f := range m.mcSaveCallback {
			f(m.EntryId)
		}
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
	cmd := []byte(c)
	sys := runtime.GOOS
	// 在windows下小于1.12是gbk编码
	if sys == constant.OS_WINDOWS && utils.CompareMcVersion(m.Version, constant.MC_LAST_UTF8_VERSION) == constant.COMPARE_LT {
		c = strings.ReplaceAll(c, "\n", "\r\n")
		cmd, _ = utils.UTF82GBK([]byte(c))
	}
	_, err := m.stdin.Write(cmd)
	return err
}

func (m *MinecraftServer) RunCommand(cmd string, params ...string) error {
	for _, param := range params {
		cmd += fmt.Sprintf(" %s", param)
	}
	return m._command(cmd)
}

// 执行tell命令
func (m *MinecraftServer) TellCommand(player string, msg string) error {
	return m.RunCommand("/tell", player, msg)
}

// 执行tellraw命令
func (m *MinecraftServer) TellrawCommand(player string, msg interface{}) error {
	rawmsg := ""
	switch msg.(type) {
	case string:
		rawmsg = utils.NewTellrowMessage().SetText(msg.(string)).JSON()
	case utils.TellrowMessage:
		rawmsg = msg.(utils.TellrowMessage).JSON()
	default:
		rawmsgbyte, _ := json.Marshal(msg)
		rawmsg = string(rawmsgbyte)
	}
	return m.RunCommand("/tellraw", player, rawmsg)
}

func (m *MinecraftServer) SayCommand(msg string) error {
	return m.RunCommand("/say", msg)
}

// validatePort
// 校验mc的端口
func (m *MinecraftServer) validatePort() (int64, error) {
	runDir := filepath.Dir(m.RunPath)
	// 新建mc配置文件
	mcConfPath := filepath.Join(runDir, constant.MC_CONF_NAME)
	if f, e := utils.CreateFile(mcConfPath); e == nil {
		f.Close()
	}
	cfg, _ := ini.Load(mcConfPath)
	var realPort int64
	if m.Port != 0 {
		realPort = m.Port
		realPortStr := strconv.FormatInt(realPort, 10)
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
	if p, e := utils.GetFreePort(realPort); p == 0 {
		if e != nil {
			m.WriteLog(e.Error(), constant.LOG_ERROR)
		}
		p, err := m.changePort(cfg, mcConfPath, 0)
		if err != nil {
			return 0, err
		}
		realPort = p
	}

	if realPort == 0 {
		return realPort, PORT_REPEAT_ERROR
	}
	return realPort, nil
}

// changePort
// 更换mc服务端端口
func (m *MinecraftServer) changePort(cfg *ini.File, path string, port int64) (int64, error) {
	// 如果可以自动更换端口就自动更换端口
	if isChange, _ := strconv.ParseBool(GetConfVal(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT)); isChange {
		unusedPort, err := utils.GetFreePort(port)
		if err != nil {
			m.WriteLog(err.Error(), constant.LOG_ERROR)
		}
		sec, _ := cfg.GetSection(ini.DefaultSection)
		unusedPortStr := strconv.FormatInt(unusedPort, 10)
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
	m.CmdStr = utils.GetCommandArr(m.Memory, m.RunPath)
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
func NewMinecraftServer(serverConf *models.ServerConf) server.MinecraftServer {
	minecraftServer := &MinecraftServer{
		ServerConf:    serverConf,
		lock:          &sync.Mutex{},
		messageChan:   make(chan string, 10),
		logger:        AddLog(serverConf.EntryId),
		ServerAdapter: &ServerAdapter{side: serverConf.Side},
	}
	RegisterCallBack(minecraftServer)
	return minecraftServer
}
