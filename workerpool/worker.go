package workerpool

type Task func() (any, error)

type Result struct {
	res any
	err error
}

type Worker struct {
	taskQueue <-chan Task
	resChan   chan<- Result
}

func (w *Worker) Run() {
	for task := range w.taskQueue {
		res, err := task()
		w.resChan <- Result{
			res: res,
			err: err,
		}
	}
}
