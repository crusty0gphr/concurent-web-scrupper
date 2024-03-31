package workerpool

import "sync"

const defaultWorkersCount = 5

type Task interface {
	Process(wg *sync.WaitGroup) // TODO: error handing here
}

type Pool struct {
	tasks        []Task
	tasksChan    chan Task
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

func (p *Pool) Run() error {
	for i := 0; i < p.workersCount; i++ {
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
	return nil // TODO: add an error return
}

func (p *Pool) worker() {
	for task := range p.tasksChan {
		task.Process(&p.wg)
	}
}
