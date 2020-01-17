package service

const (
	LOG_DEBUG   = "debug"
	LOG_INFO    = "info"
	LOG_ERROR   = "error"
	LOG_WARNING = "warning"
	LOG_FATAL   = "fatal"
)

var logLevel map[string]int

func init() {
	logLevel = make(map[string]int)
	logLevel[LOG_DEBUG] = 1
	logLevel[LOG_INFO] = 2
	logLevel[LOG_WARNING] = 3
	logLevel[LOG_ERROR] = 4
	logLevel[LOG_FATAL] = 5
}

type Log struct {
	Name      string
	Path      string
	Id        int
	level     int
	writeChan chan string
	showCodeLine bool
}

func (l *Log) Write(string, string) error {

	return nil
}

func (l *Log) SetLogLevel(level string) {
	if le, ok := logLevel[level]; ok {
		l.level = le
	}
}

func (l *Log) IsShowCodeLine(isShow bool) {
	l.showCodeLine = isShow
}

func (l *Log) GetLines(string, int, int) {

}

func (l *Log) CompressLogs(path string) {

}
