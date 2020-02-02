package _interface

type CallBack interface {
	// ChangeConfCallBack
	// 配置更改回调
	ChangeConfCallBack()

	// DestructCallBack
	// 程序退出回调
	DestructCallBack()

	// InitCallBack
	// 初始化回调
	InitCallBack()
}
