package models

import "github.com/TISUnion/most-simple-mcd/grpc/api"

// 服务器配置
type ServerConf struct {
	// EntryId
	// 实例唯一id
	EntryId string `json:"id"`

	// Name
	// 服务器名称
	Name string `json:"name"`

	// CmdStr
	// 执行的完整命令
	// 下标为0： 命令名称
	// 大于0为命令参数
	CmdStr []string `json:"cmd_str"`

	// Port
	// 启动服务器端口
	Port int64 `json:"port"`

	// RunPath
	// 运行所在工作区间
	RunPath string `json:"run_rath"`

	// IsMirror
	// 是否是镜像服务器
	IsMirror bool `json:"is_mirror"`

	// IsStartMonitor
	// 是否启动资源监听器
	IsStartMonitor bool `json:"is_start_monitor"`

	// Memory
	// 使用内存大小，单位M
	Memory int64 `json:"memory"`

	// Version
	// 服务端版本
	Version string `json:"version"`

	// GameType
	// 服务器模式
	GameType string `json:"game_type"`

	// State
	// 启动状态：0.未启动 1.启动  -1.正在启动 -2.正在关闭
	State int64 `json:"state"`

	// 本机的ip
	Ips []string `json:"ips"`
}

// 模型转数组pb对象数组
func ServerConfObjs2ServerConfPbs(srvConfList []*ServerConf) []*api.ListResp_Record {
	res := make([]*api.ListResp_Record, 0)
	for _, c := range srvConfList {
		res = append(res, &api.ListResp_Record{
			EntryId:        c.EntryId,
			Name:           c.Name,
			CmdStr:         c.CmdStr,
			Port:           c.Port,
			RunPath:        c.RunPath,
			IsMirror:       c.IsMirror,
			IsStartMonitor: c.IsStartMonitor,
			Memory:         c.Memory,
			Version:        c.Version,
			GameType:       c.GameType,
			State:          c.State,
			Ips:            c.Ips,
		})
	}
	return res
}

// 模型转pb对象
func (c *ServerConf) ServerConfObj2ServerConfPb() *api.ListResp_Record {
	return &api.ListResp_Record{
		EntryId:        c.EntryId,
		Name:           c.Name,
		CmdStr:         c.CmdStr,
		Port:           c.Port,
		RunPath:        c.RunPath,
		IsMirror:       c.IsMirror,
		IsStartMonitor: c.IsStartMonitor,
		Memory:         c.Memory,
		Version:        c.Version,
		GameType:       c.GameType,
		State:          c.State,
		Ips:            c.Ips,
	}
}
