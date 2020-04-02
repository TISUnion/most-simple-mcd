package json_struct

// 用户信息
type AdminUser struct {
	Nickname string   `json:"nickname"`
	Account  string   `json:"account"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
	Avatar   string   `json:"avatar"`
}

// 登陆token
type UserToken struct {
	Token string `json:"token"`
}

// 执行操作返回结果
type OperateResult struct {
	Status uint8 `json:"status"`
}

// 修改配置结构
type Config struct {
	ConfVal string `json:"config_val"`
	ConfKey string `json:"config_key"`
}

// 配置信息
type ConfParam struct {
	ConfVal        string `json:"config_val"`
	DefaultConfVal string `json:"-"`
	Name           string `json:"config_key"`
	Level          int    `json:"level"`
	Description    string `json:"description"`
	IsAlterable    bool   `json:"is_alterable"`
}

// 运行命令
type Command struct {
	Command string `json:"command"`

	// 1：插件运行命令  2：服务端运行命令   3：插件、服务端都运行
	Type int `json:"type"`
}

// 运行一条命令
type SingleCommand struct {
	Command  string `json:"command"`
	ServerId string `json:"id"`
	// 1：插件运行命令  2：服务端运行命令   3：插件、服务端都运行
	Type int `json:"type"`
}

type MonitorMessage struct {
	// id若为0，则为全局资源监控
	Id string `json:"id"`

	// cpu使用百分比
	CpuUsedPercent float64 `json:"cpu_used_percent"`

	// 内存使用百分比
	MemoryUsedPercent float32 `json:"memory_used_percent"`

	// 物理内存使用百分比
	VirtualMemoryUsedPercent float64 `json:"virtual_memory_used_percent"`

	// 内存使用量
	MemoryUsed uint64 `json:"memory_used"`

	// 虚拟内存使用量
	VirtualMemoryUsed uint64 `json:"virtual_memory_used"`
}