package service

import (
	"bufio"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/utils"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const (

)

var LogLevel map[string]int

func LogInit() {
	LogLevel = make(map[string]int)
	LogLevel[constant.LOG_DEBUG] = 1
	LogLevel[constant.LOG_INFO] = 2
	LogLevel[constant.LOG_WARNING] = 3
	LogLevel[constant.LOG_ERROR] = 4
	LogLevel[constant.LOG_FATAL] = 5
}

type Log struct {
	Name         string
	Path         string
	Id           int
	Level        int
	WriteChan    chan *_interface.LogMsgType
	ShowCodeLine bool
	FileObj      *os.File
	lock         *sync.Mutex // 文件锁：用于防止重新载入*os.file时，正在写入日志
}

func (l *Log) Write(logMsg *_interface.LogMsgType) {
	if l.Level > LogLevel[logMsg.Level] {
		return
	}
	l.WriteChan <- logMsg
}

func (l *Log) writeToFile() {
	select {
	case msg := <-l.WriteChan:
		if l.Level <= LogLevel[msg.Level] {
			l.lock.Lock()
			if _, err := l.FileObj.WriteString(l.getLogMsg(msg.Level, msg.Message)); err != nil {
				utils.PanicError(constant.WRITE_LOG_FAILED, err)
			}
			l.lock.Unlock()
		}
	}
}

func (l *Log) SetLogLevel(level string) {
	if le, ok := LogLevel[level]; ok {
		l.Level = le
	}
}

func (l *Log) IsShowCodeLine(isShow bool) {
	l.ShowCodeLine = isShow
}

func (l *Log) GetLines(page int, pageSize int) []string {
	result := make([]string, 0)
	scanner := bufio.NewScanner(l.FileObj)
	start := (page - 1) * pageSize
	end := page * pageSize
	n := 0
	for scanner.Scan() {
		n++
		if n > start && n < end {
			result = append(result, scanner.Text())
		}
	}
	return result
}

func (l *Log) CompressLogs(tpath string) error {
	if tpath == "" {
		extname := path.Ext(l.Path)
		dir, oldFilename := filepath.Split(l.Path)
		newFilename := strings.ReplaceAll(oldFilename, extname, ".zip")
		tpath = dir + newFilename
	}
	if err := utils.CompressFile(l.Path, tpath); err != nil {
		return err
	}
	return nil
}

func (l *Log) getLogMsg(level string, msg string) string {
	timeStr := time.Now().Format(constant.LOG_TIME_FORMAT)
	logMsg := fmt.Sprintf(constant.LOG_FORMAT, timeStr, level, msg)
	// 展示调用代码位置
	if l.ShowCodeLine {
		lineMsg := ""
		if _, file, line, ok := runtime.Caller(0); ok {
			lineMsg = fmt.Sprintf("%s:%d", file, line)
		}
		logMsg = fmt.Sprintf(constant.LOG_CODELINE_FORMAT, timeStr, level, lineMsg, msg)
	}
	return logMsg
}

func (l *Log) InitFileObj() {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.FileObj != nil {
		l.FileObj.Close()
	}
	var err error
	// 文件不存在则新建
	if l.FileObj, err = utils.CreateFile(l.Path); err != nil {
		utils.PanicError(constant.CREATE_LOG_FAILED, err)
	}
}

func (l *Log) Init() error {
	l.lock = &sync.Mutex{}
	l.InitFileObj()
	go l.writeToFile()
	return nil
}
