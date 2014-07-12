package gift

import (
	"runtime"
	"sync"
	"sync/atomic"
)

func parallelize(enabled bool, datamin, datamax int, fn func(pmin, pmax int)) {
	datasize := datamax - datamin
	partsize := datasize

	numGoroutines := 1
	if enabled {
		numProcs := runtime.GOMAXPROCS(0)
		if numProcs > 1 {
			numGoroutines = numProcs
			partsize = partsize / (numGoroutines * 10)
			if partsize < 1 {
				partsize = 1
			}
		}
	}

	if numGoroutines == 1 {
		fn(datamin, datamax)
	} else {
		var wg sync.WaitGroup
		wg.Add(numGoroutines)
		idx := int64(datamin)

		for p := 0; p < numGoroutines; p++ {
			go func() {
				defer wg.Done()
				for {
					pmin := int(atomic.AddInt64(&idx, int64(partsize))) - partsize
					if pmin >= datamax {
						break
					}
					pmax := pmin + partsize
					if pmax > datamax {
						pmax = datamax
					}
					fn(pmin, pmax)
				}
			}()
		}

		wg.Wait()
	}
}
