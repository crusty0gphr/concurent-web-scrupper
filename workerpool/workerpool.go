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

type Task func(wg *sync.WaitGroup) error

type Pool struct {
	tasksChan    chan Task
	tasks        []Task
	workersCount int
	wg           sync.WaitGroup
}

func NewWorkerPool(ops ...Option) *Pool {
	wp := &Pool{
		tasksChan:    make(chan Task),
		workersCount: defaultWorkersCount,
		wg:           sync.WaitGroup{},
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

	for i := 0; i <= p.workersCount; i++ {
		// spawn workers
		go p.worker()
	}

	p.wg.Add(len(p.tasks))
	for _, task := range p.tasks {
		p.tasksChan <- task
	}

	// all workers return
	close(p.tasksChan)

	p.wg.Wait()
	return nil
}

func (p *Pool) worker() {
	for task := range p.tasksChan {
		_ = task(&p.wg)
	}
}
