package container

type JobContainer interface {
	// RegisterJob
	// 注册定时任务
	// 第一个string为名称， 第二个为间隔具体使用查看: github.com/robfig/cron， func执行任务
	RegisterJob(string, string, func())

	StartJob(string) error
	StartJobs(...string) map[string]error
	StartAll() map[string]error

	StopJob(string)
	StopJobs(...string)

	HasJob(string) bool
}
