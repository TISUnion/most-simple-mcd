// _interface
// 接口库
package _interface

type LogMsgType struct {
	Message string
	Level string
}

// log
// 日志接口
type Log interface {

	// Write
	// 写入日志
	// 第一个string为日志等级分为：debug、info、warn、error、fatal，依次递增
	// 第二个string为写入日志内容，无需加入日志格式
	Write(*LogMsgType)

	// SetLogLevel
	// 修改日志等级, 如果日志等级比传入的等级低则不会写入该日志
	SetLogLevel(string)

	// IsShowCodeLine
	// 是否显示调用代码行数和文件
	IsShowCodeLine(bool)

	// GetLines
	// 按行分页获取日志， 第一个int为页码，第二个为一页的有多少行
	GetLines(int, int)[]string

	// CompressLogs
	// 压缩日志 传空字符串就表是压缩在当前文件夹内
	CompressLogs(path string) error
}
