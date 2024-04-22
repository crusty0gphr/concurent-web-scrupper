package workerpool

import (
	"errors"
	"fmt"
	"math"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errFailingProcess = errors.New("failing process")

func fakeProcess(i int) (int, error) {
	return i + 1, nil
}
func failingProcess(i int) (int, error) {
	if i%2 == 0 {
		return i, nil
	}
	return -1, errFailingProcess
}

func genTasks(n int, p func(i int) (int, error)) []Task {
	tasks := make([]Task, n)
	for i := 0; i < n; i++ {
		tasks[i] = func(wg *sync.WaitGroup) Error {
			defer wg.Done()
			_, err := p(i)
			return Error{
				"num": fmt.Sprintf("%d", i),
				"err": err,
			}
		}
	}
	return tasks
}

func TestErrTypes(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"ErrZeroTasks": func(t *testing.T) {
			worker := NewWorkerPool()
			err := worker.Run()

			assert.ErrorIs(t, err, ErrZeroTasks)
		},
		"ErrZeroWorkers": func(t *testing.T) {
			var zeroWorkers int
			worker := NewWorkerPool(
				WithTasks(genTasks(math.MaxUint8, fakeProcess)),
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

func TestErrorChanMessages(t *testing.T) {
	t.Run("FailingProcess", func(t *testing.T) {
		tasks := genTasks(math.MaxUint8, failingProcess)
		worker := NewWorkerPool(
			WithWorkersCount(1),
			WithTasks(tasks),
		)
		if err := worker.Run(); err != nil {
			t.Errorf("failed to run worker: %v", err)
		}

		for _, report := range worker.errSlice {
			if report["err"] != nil {
				assert.ErrorIs(t, report["err"].(error), errFailingProcess)
			}
		}
	})
}

func BenchmarkNewWorkerPool(b *testing.B) {
	printMemStats()

	tasks := make([]Task, math.MaxInt16)
	for i := 0; i < math.MaxInt16; i++ {
		tasks[i] = func(wg *sync.WaitGroup) Error {
			defer wg.Done()
			_, err := fakeProcess(i)
			return Error{
				"num": fmt.Sprintf("%d", i),
				"err": fmt.Sprintf("%v", err),
			}
		}
	}

	for i := 0; i < b.N; i++ {
		worker := NewWorkerPool(
			WithTasks(tasks),
			WithWorkersCount(25),
		)
		_ = worker.Run()
	}
	printMemStats()
	fmt.Println()
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	fmt.Printf("HeapAlloc = %.2f", float64(m.HeapAlloc)*0.000001)
	fmt.Printf("\t\tHeapObjects = %v", m.HeapObjects)
	fmt.Printf("\t\tHeapSys = %v", m.Sys)
	fmt.Printf("\t\tNumGC = %v\n", m.NumGC)
}
