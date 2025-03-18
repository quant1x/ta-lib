package algorithms

import (
	"runtime"
	"slices"
	"sync"
)

// 并行计算
// 对超大数据集可结合 goroutine 分块处理
func parallelMinMax(data []float64) (min, max float64) {
	chunkSize := len(data) / runtime.NumCPU()
	var wg sync.WaitGroup
	var mu sync.Mutex
	dataLen := len(data)
	for i := 0; i < len(data); i += chunkSize {
		wg.Add(1)
		start := i
		end := start + chunkSize
		if i+chunkSize > dataLen {
			end = dataLen
		}
		go func(chunk []float64) {
			localMin := slices.Min(chunk)
			localMax := slices.Max(chunk)
			mu.Lock()
			if localMin < min {
				min = localMin
			}
			if localMax > max {
				max = localMax
			}
			mu.Unlock()
			wg.Done()
		}(data[start:end])
	}
	wg.Wait()
	return
}
