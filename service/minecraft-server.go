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
}

func (m *MinecraftServer) ChangeConfCallBack() {
}

func (m *MinecraftServer) InitCallBack() {
}

func (m *MinecraftServer) DestructCallBack() {
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

func (m *MinecraftServer) Start() {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.isStart {
		GetLogContainerInstance().WriteLog(fmt.Sprintf("服务器: %s,重复启动", m.Name), constant.LOG_WARNING)
		return
	}
	if err := m.runProcess(); err != nil {
		return
	}
	m.isStart = true
	// TODO 加载插件
	return
}

func (m *MinecraftServer) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if !m.isStart {
		GetLogContainerInstance().WriteLog(fmt.Sprintf("服务器: %s,重复关闭", m.Name), constant.LOG_WARNING)
		return nil
	}
	m.isStart = false
	if err := m.Command("/stop"); err != nil {
		// windows下还是无法杀死进程，TODO 后期优化
		_ = m.CmdObj.Process.Kill()
	}
	return nil
}

func (m *MinecraftServer) Restart() {
	return
}

func (m *MinecraftServer) Command(c string) error {
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
	cmdObj.Dir = serverConf.RunPath
	minecraftServer := &MinecraftServer{
		ServerConf: serverConf,
		CmdObj:     cmdObj,
		stdin:      stdin,
		stdout:     stdout,
		lock:       &sync.Mutex{},
		isStart:    false,
	}
	return minecraftServer
}
