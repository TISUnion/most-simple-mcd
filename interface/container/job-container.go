package container

type JobContainer interface {
	// RegisterJob
	// 注册定时任务string为名称， uint64为间隔执行（单位：毫秒）， func执行任务
	RegisterJob(string, string, func())

	StartJob(string) error
	StartJobs(...string) map[string]error
	StartAll() map[string]error

	StopJob(string)
	StopJobs(...string)
	StopAll()
}
