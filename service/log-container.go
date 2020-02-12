package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/utils"
	"path/filepath"
	"sync"
	"time"
)

var (
	_logContainer container.LogContainer
	LogLevel      map[string]int
)

type LogContainer struct {
	NameIdMapping   map[string]int
	Logs            map[int]*Log
	LogDir          string
	lock            *sync.Mutex
	logSaveInterval string
}

func (l *LogContainer) GetLogByName(name string) _interface.Log {

	if id, ok := l.NameIdMapping[name]; ok {
		if log, ok := l.Logs[id]; ok {
			return log
		}
	}
	return nil
}

func (l *LogContainer) GetLogById(id int) _interface.Log {
	return l.Logs[id]
}

// 添加日志实例
func (l *LogContainer) AddLog(name string, params ...string) _interface.Log {
	l.lock.Lock()
	defer l.lock.Unlock()
	var (
		logLevel = constant.LOG_INFO
		path     string
		dirPath  = l.LogDir
	)
	// dirPath 为写入日志目录
	// 优先使用传入的目录，再使用容器配置目录，最后使用默认目录
	if dirPath == "" {
		dirPath = filepath.Join(GetConfVal(constant.WORKSPACE), "logs")
	}
	if len(params) > 0 {
		logLevel = params[0]
	}
	if len(params) > 1 {
		dirPath = filepath.Join(params[1])
	}
	dirPath = filepath.Join(dirPath, name)
	path = fmt.Sprintf("%s/%s.log", dirPath, time.Now().Format("2006-01-02"))
	id := GetIncreateId()
	logItem := &Log{
		Name:      name,
		Path:      path,
		Id:        id,
		Level:     LogLevel[logLevel],
		WriteChan: make(chan *_interface.LogMsgType, 10),
	}

	if err := logItem.Init(); err != nil {
		utils.PanicError(constant.CREATE_LOG_FAILED, err)
	}
	l.Logs[id] = logItem
	l.NameIdMapping[name] = id
	return logItem
}

// 每日新建新日志文件，
func (l *LogContainer) AddLogJob() {
	for _, k := range l.Logs {
		logDir := filepath.Dir(k.Path)
		logPath := fmt.Sprintf("%s/%s.log", logDir, time.Now().Format("2006-01-02"))
		if fileObj, err := utils.CreateFile(logPath); err != nil {
			utils.PanicError(constant.CREATE_LOG_FAILED, err)
		} else {
			_ = k.CompressLogs("")
			fileObj.Close()
			k.Path = logPath
			// 重载file对象，调函数是为了加锁
			k.InitFileObj()
		}

	}
}

// 把日志写到不用管道中
func (l *LogContainer) WriteLogOnChannels(msg string, level string, channels []string) {
	channels = append(channels, constant.DEFAULT_LOG_NAME)
	channels = utils.RemoveRepeatedElement(channels)
	for _, v := range channels {
		if log := l.GetLogByName(v); log != nil {
			log.Write(&_interface.LogMsgType{
				Message: msg,
				Level:   level,
			})
		}
	}
}

// 写入默认日志
func (l *LogContainer) WriteLog(params ...string) {
	log := l.GetLogByName(constant.DEFAULT_LOG_NAME)
	var (
		msg   string
		level = constant.LOG_INFO
	)
	if len(params) > 0 {
		msg = params[0]
	} else {
		return
	}

	if len(params) > 1 {
		level = params[1]
	}
	log.Write(&_interface.LogMsgType{
		Message: msg,
		Level:   level,
	})
}

// 配置修改回调
func (l *LogContainer) ChangeConfCallBack() {
	jobContainer := GetJobContainerInstance()
	logSaveInterval := GetConfVal(constant.LOG_SAVE_INTERVAL)
	if logSaveInterval != l.logSaveInterval {
		l.logSaveInterval = logSaveInterval
		jobContainer.StopJob(constant.EVERYDAY_JOB_NAME)
		// 重新注册定时清理日志任务
		jobContainer.RegisterJob(constant.EVERYDAY_JOB_NAME, logSaveInterval, l.AddLogJob)
		_ = jobContainer.StartJob(constant.EVERYDAY_JOB_NAME)
	}
}

func (l *LogContainer) DestructCallBack() {
	for _, log := range l.Logs {
		log.DestructCallBack()
	}
}

func (l *LogContainer) InitCallBack() {
	// 初始化日志map
	LogLevel = make(map[string]int)
	LogLevel[constant.LOG_DEBUG] = 1
	LogLevel[constant.LOG_INFO] = 2
	LogLevel[constant.LOG_WARNING] = 3
	LogLevel[constant.LOG_ERROR] = 4
	LogLevel[constant.LOG_FATAL] = 5

	jobContainer := GetJobContainerInstance()
	// 创建默认日志
	l.AddLog(constant.DEFAULT_LOG_NAME, constant.LOG_INFO)

	// 初始化定时清理日志任务
	jobContainer.RegisterJob(constant.EVERYDAY_JOB_NAME, GetConfVal(constant.LOG_SAVE_INTERVAL), l.AddLogJob)
	_ = jobContainer.StartJob(constant.EVERYDAY_JOB_NAME)
}

func GetLogContainerInstance() container.LogContainer {
	if _logContainer != nil {
		return _logContainer
	}
	_logContainerObj := &LogContainer{
		NameIdMapping:   make(map[string]int),
		Logs:            make(map[int]*Log),
		LogDir:          GetConfVal(constant.LOG_PATH),
		lock:            &sync.Mutex{},
		logSaveInterval: GetConfVal(constant.LOG_SAVE_INTERVAL),
	}
	// 注册回调
	RegisterCallBack(_logContainerObj)

	_logContainer = _logContainerObj
	return _logContainerObj
}

// 写入默认日志帮助函数
func WriteLogToDefault(params ...string) {
	GetLogContainerInstance().WriteLog(params...)
}

// 写入默认自定义管道日志帮助函数
func WriteLogToChannels(msg string, level string, channels []string) {
	GetLogContainerInstance().WriteLogOnChannels(msg, level, channels)
}