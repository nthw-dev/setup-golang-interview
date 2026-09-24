package counter

import (
	"sync"
	"testing"
)

func TestCounterConcurrent(t *testing.T) {
	c := NewCounter()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				c.Inc("hits")
				if i%2 == 0 {
					c.Inc("even")
				}
				_ = c.Value("hits")
			}
		}()
	}
	wg.Wait()

	if got := c.Value("hits"); got != 10000 {
		t.Errorf("hits = %d; want 10000", got)
	}
	if got := c.Value("even"); got != 5000 {
		t.Errorf("even = %d; want 5000", got)
	}
	if got := c.Value("missing"); got != 0 {
		t.Errorf("missing = %d; want 0", got)
	}
}

func TestSnapshotIsCopy(t *testing.T) {
	c := NewCounter()
	c.Inc("a")
	snap := c.Snapshot()
	snap["a"] = 999
	if c.Value("a") != 1 {
		t.Errorf("modifying snapshot changed counter: a=%d", c.Value("a"))
	}
}
