package gift

import (
	"runtime"
	"testing"
)

func testParallelizeN(enabled bool, n, procs int) bool {
	data := make([]bool, n)
	runtime.GOMAXPROCS(procs)
	parallelize(enabled, 0, n, func(start, end int) {
		for i := start; i < end; i++ {
			data[i] = true
		}
	})
	for i := 0; i < n; i++ {
		if data[i] != true {
			return false
		}
	}
	return true
}

func TestParallelize(t *testing.T) {
	for _, e := range []bool{true, false} {
		for _, n := range []int{1, 10, 100, 1000} {
			for _, p := range []int{1, 2, 4, 8, 16, 100} {
				if testParallelizeN(e, n, p) != true {
					t.Errorf("failed testParallelizeN(%v, %d, %d)", e, n, p)
				}
			}
		}
	}
}
