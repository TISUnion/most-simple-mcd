package container

import (
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/robfig/cron/v3"
	"sync"
)

type job struct {
	EntryId  cron.EntryID
	Handle   func()
	Interval string
	IsStop   bool
}

// JobContainer
// 定时任务管理器
type JobContainer struct {
	cron     *cron.Cron
	jobs     map[string]*job
	lock     *sync.Mutex
	jobNames []string
}

var JobContainerObj *JobContainer

func (jc *JobContainer) RegisterJob(name string, interval string, handle func()) {
	jc.lock.Lock()
	defer jc.lock.Unlock()
	// 已经存在任务则覆盖
	if tmpjob, ok := jc.jobs[name]; ok {
		if tmpjob.EntryId != 0 {
			jc.StopJob(name)
		}
	}

	jc.jobs[name] = &job{
		Handle:   handle,
		Interval: interval,
		IsStop:   true,
	}

	jc.jobNames = append(jc.jobNames, name)
}

func (jc *JobContainer) StartJob(name string) error {
	jc.lock.Lock()
	defer jc.lock.Unlock()
	tjob, ok := jc.jobs[name]
	if ok && tjob.IsStop {
		if id, err := jc.cron.AddFunc(tjob.Interval, tjob.Handle); err != nil {
			return err
		} else {
			// 设置开始任务的参数
			tjob.EntryId = id
			tjob.IsStop = true
		}
	}
	return nil
}

func (jc *JobContainer) StartJobs(names ...string) map[string]error {
	result := make(map[string]error)
	for _, name := range names {
		if err := jc.StartJob(name); err != nil {
			result[name] = err
		}
	}
	return result
}

func (jc *JobContainer) StartAll() map[string]error {
	return jc.StartJobs(jc.jobNames...)
}

func (jc *JobContainer) StopJob(name string) {
	jc.lock.Lock()
	defer jc.lock.Unlock()
	tjob, ok := jc.jobs[name]
	if !ok || !tjob.IsStop {
		return
	}
	if tjob.EntryId != 0 {
		jc.cron.Remove(tjob.EntryId)
	}
}

func (jc *JobContainer) StopJobs(names ...string) {
	for _, name := range names {
		jc.StopJob(name)
	}
}

func GetJobContainerObj() container.JobContainer {
	if JobContainerObj != nil {
		return JobContainerObj
	}
	jcron := cron.New()
	jcron.Start()
	JobContainerObj = &JobContainer{
		cron:     jcron,
		jobs:     make(map[string]*job),
		lock:     &sync.Mutex{},
		jobNames: make([]string, 0),
	}
	return JobContainerObj
}
