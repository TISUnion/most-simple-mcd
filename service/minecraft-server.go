package service

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"io"
	"os/exec"
	"path/filepath"
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

	// isStart
	// 是否启动
	isStart bool

	// messageChan
	// 玩家发言存储chan
	messageChan chan *server.ReciveMessageType

	// EntryId
	// 实例唯一id
	entryId int
}

func (m *MinecraftServer) GetServerConf() *json_struct.ServerConf {
	return m.ServerConf
}

func (m *MinecraftServer) SetMaxMinMemory(max int, min int) {
	if max > 0 {
		m.MaxMemory = max
	}

	if min > 0 {
		m.MinMemory = min
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
}

func (m *MinecraftServer) DestructCallBack() {
	_ = m.Stop()
	_ = m.stdin.Close()
	_ = m.stdout.Close()
}

// 启动进程
func (m *MinecraftServer) runProcess() error {
	// 校验eula
	if err := m.validateEula(); err != nil {
		return err
	}
	if port, err := m.validatePort(); err != nil {
		m.Port = port
		return err
	}
	if err := m.CmdObj.Start(); err != nil {
		return err
	}
	m.Pid = m.CmdObj.Process.Pid
	return nil
}

func (m *MinecraftServer) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.isStart {
		GetLogContainerInstance().WriteLog(fmt.Sprintf("服务器: %s,重复启动", m.Name), constant.LOG_WARNING)
		return nil
	}
	if err := m.runProcess(); err != nil {
		return err
	}
	m.isStart = true
	// TODO 加载插件
	return nil
}

func (m *MinecraftServer) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.isStart {
		GetLogContainerInstance().WriteLog(fmt.Sprintf("服务器: %s,重复关闭", m.Name), constant.LOG_WARNING)
		return nil
	}
	m.isStart = false
	if err := m._command("/stop"); err != nil {
		// windows下还是无法杀死进程，TODO 后期优化
		_ = m.CmdObj.Process.Kill()
	}

	// 重置cmd对象
	m.resetCmdObj()
	return nil
}

func (m *MinecraftServer) Restart() error {
	if m.isStart {
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
		GetLogContainerInstance().WriteLog(errMsg)
		return []byte{}, errors.New(errMsg)
	}
	// 如果一次的数据为1024，就多次获取
	if n == MAX_SIZE {
		for {
			subBuff := make([]byte, MAX_SIZE)
			subN, subErr := m.stdout.Read(buff)
			if subErr != nil {
				errMsg := fmt.Sprintf("服务器: %s，已关闭。因为%v", m.Name, err)
				GetLogContainerInstance().WriteLog(errMsg)
				return []byte{}, errors.New(errMsg)
			}
			subBuff = subBuff[:subN]
			buff = append(buff, subBuff...)
			if subN != MAX_SIZE {
				break
			}
		}
	}
	return buff, nil
}

// 获取消息，并写入到管道中
func (m *MinecraftServer) reciveMessageToChan() {
	for {
		everyBuff, err := m.resiveOneMessage()
		if err != nil {
			return
		}
		m.messageChan <- &server.ReciveMessageType{
			OriginData: everyBuff,
			ServerId:   m.entryId,
		}
	}
}

// TODO 处理消息
func (m *MinecraftServer) handleMessage() {
	for {
		msg := <-m.messageChan
		// TODO 分发给各插件

		fmt.Print(string(msg.OriginData))
	}
}

func (m *MinecraftServer) Command(c string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._command(c)
}

func (m *MinecraftServer) _command(c string) error {
	_, err := m.stdin.Write([]byte(c))
	return err
}

// validatePort
// 校验mc的端口
func (m *MinecraftServer) validatePort() (int, error) {
	mcConfPath := filepath.Join(m.RunPath, constant.MC_CONF_NAME)
	if f, e := utils.CreateFile(mcConfPath); e == nil {
		f.Close()
	}
	cfg, err := ini.Load(mcConfPath)
	var realPort int
	// 没有配置文件或者配置不完整
	if err != nil || !cfg.Section("").HasKey(constant.PORT_TEXT) {
		realPort = constant.DEFAULT_PORT
	} else {
		realPort, _ = cfg.Section("").Key(constant.PORT_TEXT).Int()
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
	if isChange, _ := strconv.ParseBool(GetConfInstance().GetConfVal(constant.IS_AUTO_CHANGE_MC_SERVER_REPEAT_PORT)); isChange {
		unusedPort, _ := utils.GetFreePort(port)
		sec, err := cfg.GetSection(ini.DefaultSection)
		if err != nil {
			return 0, err
		}
		unusedPortStr := strconv.Itoa(unusedPort)
		// 重新配置文件
		if sec.HasKey(constant.PORT_TEXT) {
			sec.Key(constant.PORT_TEXT).SetValue(unusedPortStr)
		} else {
			_, _ = sec.NewKey(constant.PORT_TEXT, unusedPortStr)
		}
		if err := cfg.SaveTo(path); err != nil {
			return 0, err
		}
		return unusedPort, nil
	} else {
		msg := fmt.Sprintf("服务端：%s，对应的服务器端口已被其他程序占用，请更换端口或者开启自动更换端口", m.Name)
		GetLogContainerInstance().WriteLog(msg)
		return 0, PORT_REPEAT_ERROR
	}
}

// validateEula
// 校验mc的eula文件
func (m *MinecraftServer) validateEula() error {
	path := filepath.Join(m.RunPath, constant.EULA_FILE_NAME)
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

func (m *MinecraftServer) resetCmdObj() {
	_ = m.stdin.Close()
	_ = m.stdout.Close()
	m.CmdObj = exec.Command(m.CmdStr[0], m.CmdStr[1:]...)
	m.stdin, _ = m.CmdObj.StdinPipe()
	m.stdout, _ = m.CmdObj.StdoutPipe()
	m.isStart = false
	m.CmdObj.Dir = m.RunPath
}

// NewMinecraftServer
// 新建一个mc服务端进程
func NewMinecraftServer(serverConf *json_struct.ServerConf) server.MinecraftServer {
	cmdObj := exec.Command(serverConf.CmdStr[0], serverConf.CmdStr[1:]...)
	stdin, err := cmdObj.StdinPipe()
	if err != nil {
		return nil
	}
	stdout, err := cmdObj.StdoutPipe()
	if err != nil {
		return nil
	}
	// 设置工作区间
	cmdObj.Dir = serverConf.RunPath
	minecraftServer := &MinecraftServer{
		ServerConf:  serverConf,
		CmdObj:      cmdObj,
		stdin:       stdin,
		stdout:      stdout,
		lock:        &sync.Mutex{},
		isStart:     false,
		messageChan: make(chan *server.ReciveMessageType, 10),
		entryId:     GetIncreateId(),
	}
	RegisterCallBack(minecraftServer)
	// 开启发送和接受消息
	go minecraftServer.reciveMessageToChan()
	go minecraftServer.handleMessage()
	return minecraftServer
}
