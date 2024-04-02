package workerpool

import (
	"math"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func fakeProcess(i int) (int, error) {
	return i + 1, nil
}

func TestNewWorkerPool(t *testing.T) {
	tasks := make([]Task, math.MaxInt8)
	for i := 0; i < math.MaxInt8; i++ {
		tasks[i] = func(wg *sync.WaitGroup) error {
			defer wg.Done()
			_, err := fakeProcess(i)
			return err
		}
	}

	tests := map[string]func(t *testing.T){
		"zero tasks": func(t *testing.T) {
			worker := NewWorkerPool()
			err := worker.Run()

			assert.ErrorIs(t, err, ErrZeroTasks)
		},
		"zero workers": func(t *testing.T) {
			var zeroWorkers int
			worker := NewWorkerPool(
				WithTasks(tasks),
				WithWorkersCount(zeroWorkers),
			)
			err := worker.Run()

			assert.ErrorIs(t, err, ErrZeroWorkers)
		},
	}

	for name, fn := range tests {
		test := fn
		t.Run(name, test)
	}
}

func BenchmarkNewWorkerPool(b *testing.B) {
	tasks := make([]Task, math.MaxInt8)
	for i := 0; i < math.MaxInt8; i++ {
		tasks[i] = func(wg *sync.WaitGroup) error {
			defer wg.Done()
			_, err := fakeProcess(i)
			return err
		}
	}

	for i := 0; i < b.N; i++ {
		worker := NewWorkerPool(
			WithTasks(tasks),
		)
		_ = worker.Run()
	}
}
