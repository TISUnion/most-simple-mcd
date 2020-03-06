package json_struct

type MonitorMessage struct {
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
