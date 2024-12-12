package parallel

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParallelProcessor(t *testing.T) {
	tests := []struct {
		name       string
		input      []int
		numWorkers int
		process    processFunc[int, int]
		wantErr    bool
		expected   []int
	}{
		{
			name:       "basic multiplication",
			input:      []int{1, 2, 3, 4, 5},
			numWorkers: 2,
			process: func(idx int, item int) (int, error) {
				return item * 2, nil
			},
			wantErr:  false,
			expected: []int{2, 4, 6, 8, 10},
		},
		{
			name:       "empty input",
			input:      []int{},
			numWorkers: 2,
			process: func(idx int, item int) (int, error) {
				return item, nil
			},
			wantErr:  false,
			expected: []int{},
		},
		{
			name:       "with error",
			input:      []int{1, 2, 3, 4, 5},
			numWorkers: 2,
			process: func(idx int, item int) (int, error) {
				if item == 3 {
					return 0, errors.New("error processing item 3")
				}
				return item * 2, nil
			},
			wantErr: true,
		},
		{
			name:       "auto worker count",
			input:      []int{1, 2, 3, 4, 5},
			numWorkers: 0,
			process: func(idx int, item int) (int, error) {
				return item * 2, nil
			},
			wantErr:  false,
			expected: []int{2, 4, 6, 8, 10},
		},
		{
			name:       "concurrent execution",
			input:      []int{1, 2, 3, 4, 5},
			numWorkers: 2,
			process: func(idx int, item int) (int, error) {
				time.Sleep(100 * time.Millisecond)
				return item * 2, nil
			},
			wantErr:  false,
			expected: []int{2, 4, 6, 8, 10},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start := time.Now()
			results, err := Processor(tt.input, tt.numWorkers, tt.process)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, results)

			if tt.name == "concurrent execution" {
				duration := time.Since(start)
				// With 2 workers, 5 items should take ~300ms (not 500ms)
				assert.Less(t, duration, 400*time.Millisecond,
					"execution took too long: %v", duration)
			}
		})
	}
}
