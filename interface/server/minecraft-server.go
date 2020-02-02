package server

type MinecraftServer interface {
	BasicServer
	Command(string, ...interface{}) error //执行命令
}
