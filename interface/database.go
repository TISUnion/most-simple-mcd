package _interface

import "time"

// Database
// 数据库, 对github.com/dgraph-io/badger包的二次封装
type Database interface {
	CallBack

	// 获取值
	Get(string) string

	// 设置值
	Set(string, string)

	// 设置值同时设置过期时间
	SetWiteTTL(string, string, time.Duration)
}
