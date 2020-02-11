package server

type ReciveMessageType struct {
	Player string
	Time string
	Speak string
	OriginData []byte
}

type MinecraftServer interface {
	BasicServer
	// 第一参数为命令，后面为参数
	Command(string) error //执行命令
}
