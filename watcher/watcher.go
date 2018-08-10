package watcher

import (
	"fmt"
	"os"

	"github.com/uphy/elastic-watcher/config"
)

type Watcher struct {
	watches      []*WatchJob
	scheduler    *Scheduler
	globalConfig *config.Config
}

type WatchJob struct {
	w        *Watch
	schedule string
}

func New(globalConfig *config.Config) *Watcher {
	return &Watcher{
		scheduler:    newScheduler(),
		globalConfig: globalConfig,
	}
}

func (a *WatchJob) Run() {
	if err := a.w.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func (w *Watcher) AddWatch(watchConfig *WatchConfig) {
	w.scheduler.AddTask(watchConfig.Trigger.Schedule, &WatchJob{
		NewWatch(w.globalConfig, watchConfig),
		*watchConfig.Trigger.Schedule.Cron,
	})
}

func (w *Watcher) Start() {
	w.scheduler.Start()
}

func (w *Watcher) Stop() {
	w.scheduler.Stop()
}
