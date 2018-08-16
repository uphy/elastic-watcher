package context

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"

	"github.com/Sirupsen/logrus"
	multierror "github.com/hashicorp/go-multierror"
)

type (
	TaskRunner struct {
		id       string
		logger   *logrus.Entry
		parent   *TaskRunner
		children []*TaskRunner
		ctx      ExecutionContext
		stopped  bool
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

func newTaskRunner(ctx ExecutionContext) *TaskRunner {
	return &TaskRunner{
		id:       ctx.ID(),
		logger:   ctx.Logger(),
		parent:   nil,
		children: nil,
		ctx:      ctx,
		stopped:  false,
	}
}

func (t *TaskRunner) Init() {
	t.children = nil
}

func NewTask(f TaskFunc) Task {
	return &task{f}
}

func (t *task) Run(ctx ExecutionContext) error {
	return t.t(ctx)
}

func (t *TaskRunner) RunFunc(f TaskFunc) error {
	return t.Run(NewTask(f))
}

func (t *TaskRunner) Run(task Task) error {
	workers := []*TaskRunner{}
	if len(t.children) == 0 {
		workers = append(workers, t)
	} else {
		for _, w := range t.children {
			workers = append(workers, w)
		}
	}
	t.logger.WithFields(logrus.Fields{
		"type":    reflect.TypeOf(task),
		"workers": len(workers),
	}).Debug("Running task...")
	t.logWorkers()
	for _, worker := range workers {
		if err := worker.run(task); err != nil {
			if err == ErrStop {
				continue
			}
			return err
		}
	}
	return nil
}

func (t *TaskRunner) remove(child *TaskRunner) error {
	for i, c := range t.children {
		if c.id == child.id {
			t.children = append(t.children[:i], t.children[i+1:]...)
			return nil
		}
	}
	return errors.New("no such child")
}

func (t *TaskRunner) stop() {
	t.stopped = true
	if t.parent != nil {
		t.parent.remove(t)
	}
}

func (t *TaskRunner) run(task Task) error {
	if t.stopped {
		return errors.New("runner stopped")
	}
	if err := task.Run(t.ctx); err != nil {
		if err == ErrStop {
			t.stop()
			t.logger.WithFields(logrus.Fields{
				"type":   reflect.TypeOf(task),
				"worker": t.id,
			}).Debug("worker stopped")
			return err
		}
		return err
	}
	splitted := consumeSplittedPayload(t.ctx)
	if splitted != nil {
		t.logger.WithFields(logrus.Fields{
			"current": len(t.children),
			"new":     len(t.children) + len(splitted) - 1,
		}).Debug("Splitting workers...")
		if err := t.forkWorker(splitted); err != nil {
			return err
		}
		t.logWorkers()
	}
	return nil
}

func (t *TaskRunner) forkWorker(splittedPayload []JSONObject) error {
	for _, p := range splittedPayload {
		c, err := wrapContext(t.ctx, false)
		if err != nil {
			return err
		}
		c.SetPayload(p)
	}
	return nil
}

func (t *TaskRunner) addWorker(ctx ExecutionContext) *TaskRunner {
	child := newTaskRunner(ctx)
	t.children = append(t.children, child)
	child.parent = t
	return child
}

func (t *TaskRunner) root() *TaskRunner {
	if t.parent == nil {
		return t
	}
	return t.parent.root()
}

func (t *TaskRunner) logWorkers() {
	if t.ctx.GlobalConfig().Debug {
		t.logger.Debugf("workers:\n%s", t.root().strace())
	}
}

func (t *TaskRunner) strace() string {
	buf := new(bytes.Buffer)
	t.trace("", buf)
	return buf.String()
}
func (t *TaskRunner) trace(indent string, w io.Writer) {
	fmt.Fprintf(w, "%s%s: %v\n", indent, t.id, t.ctx.Payload())
	if t.children != nil {
		for _, c := range t.children {
			c.trace(indent+"  ", w)
		}
	}
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
			wrappedCtx, err := wrapContext(ctx, true)
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
