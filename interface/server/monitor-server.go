package server

type MonitorMessage struct {
	// cpu使用百分比
	CpuUsedPercent float64
	// 内存使用百分比
	MemoryUsedPercent float32
	// 物理内存使用百分比
	VirtualMemoryUsedPercent float64

	// 内存使用量
	MemoryUsed uint64
	// 虚拟内存使用量
	VirtualMemoryUsed uint64
}

type MonitorServer interface {
	BasicServer

	GetMessageChan() chan *MonitorMessage
}
