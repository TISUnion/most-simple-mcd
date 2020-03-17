package modules

import (
	"bufio"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/utils"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Log struct {
	Name      string
	Path      string
	Id        int
	Level     int
	WriteChan chan *_interface.LogMsgType
	FileObj   *os.File
	lock      *sync.Mutex // 文件锁：用于防止重新载入*os.file时，正在写入日志
}

// 直接写入日志，不进行格式化
func (l *Log) Write(p []byte) (n int, err error) {
	l.WriteLog(&_interface.LogMsgType{
		Message:     string(p),
		Level:       constant.LOG_INFO,
		IsNotFormat: true,
	})
	return 0, nil
}

func (l *Log) Errorf(f string, v ...interface{}) {
	l.WriteLog(&_interface.LogMsgType{
		Message: fmt.Sprintf(f, v),
		Level:   constant.LOG_ERROR,
	})
}

func (l *Log) Warningf(f string, v ...interface{}) {
	l.WriteLog(&_interface.LogMsgType{
		Message: fmt.Sprintf(f, v),
		Level:   constant.LOG_WARNING,
	})
}

func (l *Log) Infof(f string, v ...interface{}) {
	l.WriteLog(&_interface.LogMsgType{
		Message: fmt.Sprintf(f, v),
		Level:   constant.LOG_INFO,
	})
}

func (l *Log) Debugf(f string, v ...interface{}) {
	l.WriteLog(&_interface.LogMsgType{
		Message: fmt.Sprintf(f, v),
		Level:   constant.LOG_DEBUG,
	})
}

func (l *Log) ChangeConfCallBack() {}

func (l *Log) DestructCallBack() {
	l.FileObj.Close()
}

func (l *Log) InitCallBack() {
}

func (l *Log) WriteLog(logMsg *_interface.LogMsgType) {
	if l.Level > LogLevel[logMsg.Level] {
		return
	}
	l.WriteChan <- logMsg
}

func (l *Log) writeToFile() {
	for {
		select {
		case msg := <-l.WriteChan:
			l.lock.Lock()
			var data string
			// 如果是否不需要格式化
			if msg.IsNotFormat {
				data = msg.Message
			} else {
				data = l.getLogMsg(msg.Level, msg.Message)
			}
			if _, err := l.FileObj.WriteString(data); err != nil {
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
