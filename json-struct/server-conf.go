package json_struct

type ServerConf struct {
	// Name
	// 服务器名称
	Name string

	// CmdStr
	// 执行的完整命令
	// 下标为0： 命令名称
	// 大于0为命令参数
	CmdStr []string

	// Port
	// 启动服务器端口
	Port int

	// RunPath
	// 运行所在工作区间
	RunPath string

	// IsMirror
	// 是否是镜像服务器
	IsMirror bool

	// IsStartMonitor
	// 是否启动资源监听器
	IsStartMonitor bool

	// Memory
	// 使用内存大小，单位M
	Memory int
}
