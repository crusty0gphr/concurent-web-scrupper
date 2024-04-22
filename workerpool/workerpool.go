package workerpool

import (
	"errors"
	"sync"
)

var (
	ErrZeroTasks   = errors.New("no tasks to run")
	ErrZeroWorkers = errors.New("zero workers provided")
)

const defaultWorkersCount = 5

type (
	Error map[string]any
	Task  func(wg *sync.WaitGroup) Error
)

var (
	wg sync.WaitGroup
	mu sync.RWMutex
)

type Pool struct {
	tasksChan    chan Task
	errSlice     []Error
	tasks        []Task
	workersCount int
}

func NewWorkerPool(ops ...Option) *Pool {
	wp := &Pool{
		tasksChan:    make(chan Task),
		errSlice:     []Error{},
		workersCount: defaultWorkersCount,
	}

	for _, opFunc := range ops {
		opFunc(wp)
	}
	return wp
}

func (p *Pool) fiascoCheck() error {
	switch {
	case len(p.tasks) == 0:
		return ErrZeroTasks
	case p.workersCount == 0:
		return ErrZeroWorkers
	default:
		return nil
	}
}

func (p *Pool) Run() error {
	if err := p.fiascoCheck(); err != nil {
		return err
	}

	for i := 1; i <= p.workersCount; i++ {
		// spawn workers
		go p.worker()
	}

	wg.Add(len(p.tasks))
	for _, task := range p.tasks {
		p.tasksChan <- task
	}
	wg.Wait()

	// all workers return
	close(p.tasksChan)
	return nil
}

func (p *Pool) worker() {
	mu.RLock()
	defer mu.RUnlock()

	for task := range p.tasksChan {
		if err := task(&wg); err != nil {
			p.errSlice = append(p.errSlice, err)
		}
	}
}
