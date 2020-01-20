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
	increateId    int         = 0
	idLock        *sync.Mutex = &sync.Mutex{}
	_logContainer container.LogContainer
)

func getIncreateId() int {
	idLock.Lock()
	defer idLock.Unlock()
	increateId++
	return increateId
}

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

func (l *LogContainer) AddLog(name string, params ...interface{}) _interface.Log {
	l.lock.Lock()
	defer l.lock.Unlock()
	var (
		logLevel       = constant.LOG_INFO
		isShowCodeLine = false
		path           string
		dirPath        = l.LogDir
	)
	// dirPath 为写入日志目录
	// 优先使用传入的目录，再使用容器配置目录，最后使用默认目录
	if dirPath == "" {
		if currentPath, err := utils.GetCurrentPath(); err == nil {
			dirPath = filepath.Join(currentPath, "logs")
		} else {
			dirPath = filepath.Join("./", "logs")
		}
	}
	if len(params) > 0 {
		logLevel = params[0].(string)
	}
	if len(params) > 1 {
		isShowCodeLine = params[1].(bool)
	}
	if len(params) > 2 {
		dirPath = filepath.Join(params[2].(string))
	}
	dirPath = filepath.Join(dirPath, name)
	path = fmt.Sprintf("%s/%s.log", dirPath, time.Now().Format("2006-01-02"))
	id := getIncreateId()
	logItem := &Log{
		Name:         name,
		Path:         path,
		Id:           id,
		Level:        LogLevel[logLevel],
		ShowCodeLine: isShowCodeLine,
		WriteChan:    make(chan *_interface.LogMsgType),
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

func (l *LogContainer) WriteLogOnChannels(msg string, level string, channels []string) {
	channels = append(channels, constant.DEFAULT_CHANNEL)
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
	log := l.GetLogByName(constant.DEFAULT_CHANNEL)
	var (
		msg   string
		level string
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
	conf := GetConfInstance()
	jobContainer := GetJobContainerInstance()
	if conf.GetConfVal(constant.LOG_SAVE_INTERVAL) != l.logSaveInterval {
		l.logSaveInterval = conf.GetConfVal(constant.LOG_SAVE_INTERVAL)
		jobContainer.StopJob(constant.EVERYDAY_JOB_NAME)
		// 重新注册定时清理日志任务
		jobContainer.RegisterJob(constant.EVERYDAY_JOB_NAME, conf.GetConfVal(constant.LOG_SAVE_INTERVAL), l.AddLogJob)
		_ = jobContainer.StartJob(constant.EVERYDAY_JOB_NAME)
	}
}

func GetLogContainerObj() container.LogContainer {
	if _logContainer != nil {
		return _logContainer
	}
	conf := GetConfInstance()
	jobContainer := GetJobContainerInstance()
	_logContainerObj := &LogContainer{
		NameIdMapping:   make(map[string]int),
		Logs:            make(map[int]*Log),
		LogDir:          conf.GetConfVal(constant.LOG_PATH),
		lock:            &sync.Mutex{},
		logSaveInterval: conf.GetConfVal(constant.LOG_SAVE_INTERVAL),
	}
	// 创建默认日志
	_logContainerObj.AddLog(constant.DEFAULT_CHANNEL)
	_logContainer = _logContainerObj
	// 初始化定时清理日志任务
	jobContainer.RegisterJob(constant.EVERYDAY_JOB_NAME, conf.GetConfVal(constant.LOG_SAVE_INTERVAL), _logContainerObj.AddLogJob)
	_ = jobContainer.StartJob(constant.EVERYDAY_JOB_NAME)
	return _logContainerObj
}
