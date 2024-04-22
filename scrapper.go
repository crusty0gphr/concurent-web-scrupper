package concurent_web_scrupper

import (
	"sync"

	"github.com/concurent-web-scrupper/parser"
	"github.com/concurent-web-scrupper/workerpool"
)

const (
	urlKey = "url"
	errKey = "error"
)

func Run(urls []string) error {
	out := make(chan []string, len(urls))
	tasks := make([]workerpool.Task, len(urls))

	for i, url := range urls {
		tasks[i] = func(wg *sync.WaitGroup) workerpool.Error {
			res, err := parser.ExtractUrls(url)
			if err != nil {
				return workerpool.Error{
					urlKey: url,
					errKey: err,
				}
			}
			out <- res
			return nil
		}
	}

	worker := workerpool.NewWorkerPool()
	if err := worker.Run(); err != nil {
		return err
	}
	return nil
}
