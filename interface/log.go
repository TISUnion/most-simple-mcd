// _interface
// 接口库
package _interface

// log
// 日志接口
type log interface {
	// append
	// 日志尾追加
	append(string) error

	// getLines
	// 获取日志， 第一个int为开始行数，第二个为偏移量
	getLines(string, int, int)

	//CompressLogs
	// 压缩日志
	CompressLogs()
}
