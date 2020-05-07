package models



// 服务端运行状态
type ServerRunState struct {
	// State
	// 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
	State int `json:"state"`
}

// 服务器接收消息
type ReciveMessage struct {
	Player     string `json:"player"`
	Time       string `json:"time"`
	Speak      string `json:"speak"`
	OriginData []byte `json:"origin_data"`
	ServerId   string `json:"server_id"`
}


// 服务器详情
type ServerDetail struct {
	ServInfo *ServerConf   `json:"server_info"`
	PlgnInfo []*PluginInfo `json:"plugin_info"`
}

// 操作服务器
type OperateServer struct {
	ServerId []string `json:"id"`
	// 操作类型：1. 启动  2. 停止  3.重启
	OperateType int `json:"operate_type"`
}

// 操作插件
type OperatePlugin struct {
	ServerId string   `json:"server_id"`
	PluginId []string `json:"plugin_id"`
	// 操作类型：1. 启动  2. 停止
	OperateType int `json:"operate_type"`
}

// 插件命令行解析对象
type PluginCommand struct {
	Command  string
	Params   []string
	ParamMap map[string]string
}
