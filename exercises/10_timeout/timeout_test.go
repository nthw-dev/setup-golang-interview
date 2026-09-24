package timeout

import (
	"context"
	"errors"
	"runtime"
	"testing"
	"time"
)

func sleepy(d time.Duration, val string) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		select {
		case <-time.After(d):
			return val, nil
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

func TestCallWithTimeoutFast(t *testing.T) {
	got, err := CallWithTimeout(context.Background(), 100*time.Millisecond, sleepy(5*time.Millisecond, "ok"))
	if err != nil || got != "ok" {
		t.Fatalf("got %q, %v; want ok, nil", got, err)
	}
}

func TestCallWithTimeoutSlow(t *testing.T) {
	start := time.Now()
	_, err := CallWithTimeout(context.Background(), 50*time.Millisecond, sleepy(time.Second, "late"))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v; want context.DeadlineExceeded", err)
	}
	if time.Since(start) > 300*time.Millisecond {
		t.Errorf("took %v; should return right after the timeout", time.Since(start))
	}
}

func TestCallWithTimeoutNoLeak(t *testing.T) {
	base := runtime.NumGoroutine()
	// fn นี้ไม่สนใจ ctx — ทำงานจนเสร็จเองใน 100ms
	stubborn := func(context.Context) (string, error) {
		time.Sleep(100 * time.Millisecond)
		return "done", nil
	}
	for range 5 {
		_, _ = CallWithTimeout(context.Background(), 10*time.Millisecond, stubborn)
	}
	deadline := time.Now().Add(time.Second)
	for runtime.NumGoroutine() > base {
		if time.Now().After(deadline) {
			t.Fatalf("goroutine leak: %d running (base %d) — hint: buffered channel",
				runtime.NumGoroutine(), base)
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func TestFirstSuccess(t *testing.T) {
	errA := errors.New("a failed")
	failFast := func(context.Context) (string, error) { return "", errA }

	start := time.Now()
	got, err := FirstSuccess(context.Background(),
		failFast,
		sleepy(30*time.Millisecond, "fast"),
		sleepy(2*time.Second, "slow"),
	)
	if err != nil || got != "fast" {
		t.Fatalf("got %q, %v; want fast, nil", got, err)
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Errorf("took %v; should not wait for slow function", time.Since(start))
	}
}

func TestFirstSuccessAllFail(t *testing.T) {
	errA, errB := errors.New("a"), errors.New("b")
	_, err := FirstSuccess(context.Background(),
		func(context.Context) (string, error) { return "", errA },
		func(context.Context) (string, error) { return "", errB },
	)
	if !errors.Is(err, errA) || !errors.Is(err, errB) {
		t.Fatalf("err = %v; want both errA and errB (errors.Join)", err)
	}
	if _, err := FirstSuccess(context.Background()); !errors.Is(err, ErrNoFuncs) {
		t.Errorf("no funcs: err = %v; want ErrNoFuncs", err)
	}
}
