package workerpool

type Pool struct {
	concurrencyLevel int
	taskQueue        chan Task
	resChan          chan Result
}

func (p Pool) Run() {
	for i := 0; i < p.concurrencyLevel; i++ {
		w := Worker{
			taskQueue: p.taskQueue,
			resChan:   p.resChan,
		}
		w.Run()
	}
}

func (p Pool) Add(task Task) {
	p.taskQueue <- task
}

func (p Pool) Read() Result {
	return <-p.resChan
}
