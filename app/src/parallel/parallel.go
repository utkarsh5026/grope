package parallel

import (
	"fmt"
	"runtime"
	"sync"
)

type Result[T any] struct {
	Index int
	Data  T
	Err   error
}

type processFunc[In any, Out any] func(itemIdx int, item In) (Out, error)

// Processor processes a slice of items in parallel using a specified number of workers.
// It takes a slice of input items of type In, processes each item using the provided
// processFunc, and returns a slice of output items of type Out.
//
// Parameters:
//   - items: The slice of input items to process
//   - numWorkers: The number of worker goroutines to use. If <= 0, uses runtime.NumCPU()
//   - processFunc: The function to process each item, taking (index, item) and returning (result, error)
//
// Returns:
//   - []Out: Slice of processed results in the same order as input items
//   - error: First error encountered during processing, if any
//
// The function maintains the order of results relative to input items. If any worker
// encounters an error, processing continues but the function will return an error
// containing all errors encountered.
func Processor[In any, Out any](items []In, numWorkers int, processFunc processFunc[In, Out]) ([]Out, error) {

	if len(items) == 0 {
		return []Out{}, nil
	}

	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	if numWorkers > len(items) {
		numWorkers = len(items)
	}

	itemChan := make(chan int)
	resultsChan := make(chan Result[Out], len(items))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for itemIdx := range itemChan {
				item := items[itemIdx]
				result, err := processFunc(itemIdx, item)
				resultsChan <- Result[Out]{
					Index: itemIdx,
					Data:  result,
					Err:   err,
				}
			}
		}()
	}

	go func() {
		for idx := range items {
			itemChan <- idx
		}
		close(itemChan)
		wg.Wait()
		close(resultsChan)
	}()

	results := make([]Out, len(items))
	var errors []error

	for result := range resultsChan {
		if result.Err != nil {
			errors = append(errors, result.Err)
		} else {
			results[result.Index] = result.Data
		}
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("errors occurred in parallel execution: %v", errors)
	}

	return results, nil
}
