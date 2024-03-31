package workerpool

import (
	"math"
	"sync"
	"testing"
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

	worker := NewWorkerPool(
		WithTasks(tasks),
	)
	_ = worker.Run()
}

func BenchmarkName(b *testing.B) {
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
