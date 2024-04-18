package workerpool

import (
	"fmt"
	"math"
	"runtime"
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
	printMemStats()

	tasks := make([]Task, math.MaxInt16)
	for i := 0; i < math.MaxInt16; i++ {
		tasks[i] = func(wg *sync.WaitGroup) error {
			defer wg.Done()
			_, err := fakeProcess(i)
			return err
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
	fmt.Printf("\t\tHeapObjects = %v", (m.HeapObjects))
	fmt.Printf("\t\tHeapSys = %v", (m.Sys))
	fmt.Printf("\t\tNumGC = %v\n", m.NumGC)
}
