package context

import (
	"reflect"
	"sync"

	"github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
)

type (
	TaskRunner struct {
		workers []*worker
		logger  *logrus.Entry
		baseCtx ExecutionContext
	}
	worker struct {
		id  string
		ctx ExecutionContext
	}
	TaskFunc func(ctx ExecutionContext) error
	Task     interface {
		Run(ctx ExecutionContext) error
	}
	task struct {
		t TaskFunc
	}
	parallelTask struct {
		tasks []Task
	}
	stop struct {
	}
)

func (s *stop) Error() string {
	return "stop"
}

var ErrStop = &stop{}

func NewTaskRunner(ctx ExecutionContext) *TaskRunner {
	r := &TaskRunner{
		logger:  ctx.Logger(),
		baseCtx: ctx,
	}
	r.addWorker(ctx)
	return r
}

func (t *TaskRunner) Init() {
	t.workers = nil
	t.addWorker(t.baseCtx)
}

func (t *TaskRunner) addWorker(ctx ExecutionContext) *worker {
	w := &worker{
		id:  generateID(),
		ctx: ctx,
	}
	t.workers = append(t.workers, w)
	return w
}

func (t *task) Run(ctx ExecutionContext) error {
	return t.t(ctx)
}

func (t *TaskRunner) RunFunc(f TaskFunc) error {
	return t.Run(&task{f})
}

func (t *TaskRunner) Run(task Task) error {
	workers := t.workers
	t.logger.WithFields(logrus.Fields{
		"type":    reflect.TypeOf(task),
		"workers": len(workers),
	}).Debug("Running task...")
	for i, w := range workers {
		if err := w.run(task); err != nil {
			if err == ErrStop {
				t.stopWorker(i)
				continue
			}
			return err
		}
		splitted := consumeSplittedPayload(w.ctx)
		if len(splitted) > 1 {
			t.logger.WithFields(logrus.Fields{
				"current": len(t.workers),
				"new":     len(t.workers) + len(splitted) - 1,
			}).Debug("Splitting workers...")
			t.stopWorker(i)
			if err := t.forkWorker(w, splitted); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *TaskRunner) stopWorker(i int) {
	t.workers = append(t.workers[:i], t.workers[i+1:]...)
}

func (t *TaskRunner) forkWorker(w *worker, splittedPayload []JSONObject) error {
	for _, p := range splittedPayload {
		c, err := Wrap(w.ctx)
		if err != nil {
			return err
		}
		c.SetPayload(p)
		t.addWorker(c)
	}
	return nil
}

func (w *worker) run(task Task) error {
	return task.Run(w.ctx)
}

func NewParallelTask(tasks []Task) Task {
	return &parallelTask{tasks}
}

func (p *parallelTask) Run(ctx ExecutionContext) error {
	var errs error
	errmutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for _, task := range p.tasks {
		wg.Add(1)
		go func(task Task, ctx ExecutionContext) {
			wrappedCtx, err := Wrap(ctx)
			if err != nil {
				errs = multierror.Append(errs, err)
				return
			}
			if err := task.Run(wrappedCtx); err != nil {
				errmutex.Lock()
				defer errmutex.Unlock()
				errs = multierror.Append(errs, err)
			}
			wg.Done()
		}(task, ctx)
	}
	wg.Wait()
	return errs
}

/*
func (t *TaskRunner) RunParallelFunc(f ...TaskFunc) error {
	tasks := []Task{}
	for _, ff := range f {
		tasks = append(tasks, &task{ff})
	}
	return t.RunParallel(tasks...)
}

func (t *TaskRunner) RunParallel(tasks ...Task) error {
	var err error
	errmutex := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for _, task := range tasks {
		wg.Add(1)
		go func(task Task) {
			if e := task.Run(Wrap(ctx)); err != nil {
				errmutex.Lock()
				defer errmutex.Unlock()
				err = multierror.Append(err, e)
			}
			wg.Done()
		}(ctx)
	}
	wg.Wait()
	return err
}
*/
