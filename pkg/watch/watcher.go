package watch

import (
	"fmt"
	"os"

	"github.com/uphy/elastic-watcher/pkg/config"
)

type Watcher struct {
	watches      []*WatchJob
	scheduler    *Scheduler
	globalConfig *config.Config
}

type WatchJob struct {
	w *Watch
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

func (w *Watcher) AddWatch(watchConfig *WatchConfig) error {
	schedules, err := watchConfig.Trigger.Schedule.CronSchedules()
	if err != nil {
		return err
	}
	for _, s := range schedules {
		if err := w.scheduler.AddTask(s, &WatchJob{NewWatch(w.globalConfig, watchConfig)}); err != nil {
			return err
		}
	}
	return nil
}

func (w *Watcher) Start() {
	w.scheduler.Start()
}

func (w *Watcher) Stop() {
	w.scheduler.Stop()
}
