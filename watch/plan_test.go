package watch

import (
	"testing"
	"time"
)

func init() {
	watchFuncFactory["noop"] = noopWatch
}

func noopWatch(params map[string][]string) (WatchFunc, error) {
	fn := func(p *WatchPlan) (uint64, interface{}, error) {
		idx := p.lastIndex + 1
		return idx, idx, nil
	}
	return fn, nil
}

func TestRun_Stop(t *testing.T) {
	plan, err := Parse("type:noop")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	var expect uint64 = 1
	plan.Handler = func(idx uint64, val interface{}) {
		if idx != expect {
			t.Fatalf("Bad: %d %d", expect, idx)
		}
		if val != expect {
			t.Fatalf("Bad: %d %d", expect, val)
		}
		expect++
	}

	time.AfterFunc(10*time.Millisecond, func() {
		plan.Stop()
	})

	err = plan.Run("127.0.0.1:8500")
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if expect == 1 {
		t.Fatalf("Bad: %d", expect)
	}
}
