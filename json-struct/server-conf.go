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

	// MaxMemory
	// 最大运行内存单位M
	MaxMemory int

	// MinMemory
	// 最小运行内存单位M
	MinMemory int

	// IsStartGui
	// 是否启用gui
	IsStartGui bool
}
