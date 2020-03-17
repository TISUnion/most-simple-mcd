package modules

import (
	"errors"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
	"sync"
)

var (
	NOT_EXIST_PROCESS = errors.New("不存在该pid的进程")
)

type MonitorServer struct {
	messageChan chan *json_struct.MonitorMessage
	serverId    string
	serverPid   int
	processObj  *process.Process
	jobName     string
	lock        *sync.Mutex
}

func (m *MonitorServer) GetMessageChan() chan *json_struct.MonitorMessage {
	return m.messageChan
}

func (m *MonitorServer) ChangeConfCallBack() {
	_ = m.Restart()
}

func (m *MonitorServer) DestructCallBack() {
	_ = m._stop()
	m.messageChan = nil
}

func (m *MonitorServer) InitCallBack() {

}

func (m *MonitorServer) Start() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._start()
}

func (m *MonitorServer) _start() error {
	pid := utils.IntToint32(m.serverPid)
	if ok, _ := process.PidExists(pid); !ok {
		return NOT_EXIST_PROCESS
	}

	if m.processObj == nil {
		m.processObj, _ = process.NewProcess(pid)
	}

	interval := GetConfVal(constant.MONITOR_INTERVAL)
	cronStr := fmt.Sprintf("@every %s", interval)
	jobC := GetJobContainerInstance()
	jobC.RegisterJob(m.jobName, cronStr, m.GetMonitorMessage)
	_ = jobC.StartJob(m.jobName)
	return nil
}

func (m *MonitorServer) GetMonitorMessage() {
	cpuPercent, _ := m.processObj.CPUPercent()
	memoryPercent, _ := m.processObj.MemoryPercent()
	memoryInfo, _ := m.processObj.MemoryInfo()

	virtualMem, _ := mem.VirtualMemory()
	msg := &json_struct.MonitorMessage{
		CpuUsedPercent:           cpuPercent,
		MemoryUsedPercent:        memoryPercent,
		VirtualMemoryUsedPercent: utils.Uint64Tofloat64(memoryInfo.VMS) / utils.Uint64Tofloat64(virtualMem.Total) * 100,
		MemoryUsed:               memoryInfo.RSS,
		VirtualMemoryUsed:        memoryInfo.VMS,
	}
	m.messageChan <- msg
}

func (m *MonitorServer) Stop() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m._stop()
}

func (m *MonitorServer) _stop() error {
	jobC := GetJobContainerInstance()
	jobC.StopJob(m.jobName)
	// 重置管道，防止有数据遗留在管道中，下次被读出
	m.messageChan = make(chan *json_struct.MonitorMessage, 10)
	return nil
}

func (m *MonitorServer) Restart() error {
	m.lock.Lock()
	defer m.lock.Unlock()
	if err := m._stop(); err != nil {
		return err
	}
	if err := m._start(); err != nil {
		return err
	}

	return nil
}

func NewMonitorServer(id string, pid int) server.MonitorServer {
	ms := &MonitorServer{
		messageChan: make(chan *json_struct.MonitorMessage, 10),
		serverId:    id,
		serverPid:   pid,
		jobName:     fmt.Sprintf("monitor:%d", id),
		lock:        &sync.Mutex{},
	}
	return ms
}
