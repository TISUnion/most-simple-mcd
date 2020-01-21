package _interface

import "time"

// Database
// 数据库, 对github.com/dgraph-io/badger包的二次封装
type Database interface {
	CallBack

	// Get
	// 获取值
	Get(string) string

	// Set
	// 设置值
	Set(string, string)

	// SetWiteTTL
	// 设置值同时设置过期时间
	SetWiteTTL(string, string, time.Duration)
}
