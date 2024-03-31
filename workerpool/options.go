package workerpool

type Option func(p *Pool)

func WithWorkersCount(c int) Option {
	return func(p *Pool) {
		p.workersCount = c
	}
}

func WithTasks(t []Task) Option {
	return func(p *Pool) {
		p.tasks = t
	}
}
