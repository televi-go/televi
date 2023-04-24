package update

import (
	"github.com/televi-go/televi/telegram"
	"github.com/televi-go/televi/telegram/bot"
	"github.com/televi-go/televi/util"
	"sync"
)

type CompareAction struct {
	CanBeUpdated  bool
	UpdateActions []Update
}

type Update interface {
	GetRequest(destination telegram.Destination) telegram.Request
	// InflateResult may result in new request
	InflateResult(response telegram.Response) Update
}

type CompareResult struct {
	Parallel    []Update
	Consecutive []Update
}

func (result CompareResult) AddPart(part CompareResult) CompareResult {
	return CompareResult{Parallel: util.Merge(result.Parallel, part.Parallel), Consecutive: util.Merge(result.Consecutive, part.Consecutive)}
}

func (result CompareResult) Run(api *bot.Api, destination telegram.Destination) error {
	wg := sync.WaitGroup{}
	wg.Add(len(result.Parallel))
	errorChannel := make(chan error, len(result.Parallel))
	for _, update := range result.Parallel {
		go func(update Update) {
			defer wg.Done()
			resp, err := api.Request(update.GetRequest(destination))
			if err != nil {
				errorChannel <- err
				return
			}
			update.InflateResult(resp)
		}(update)
	}
	for _, update := range result.Consecutive {
		resp, err := api.Request(update.GetRequest(destination))
		if err != nil {
			return err
		}
		update.InflateResult(resp)
	}
	wg.Wait()
	return nil
}
