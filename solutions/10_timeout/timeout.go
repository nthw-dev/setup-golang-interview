// Package timeout — เฉลยข้อ 10: context.WithTimeout, select, buffered channel กัน leak
package timeout

import (
	"context"
	"errors"
	"time"
)

var ErrNoFuncs = errors.New("no functions given")

type result struct {
	val string
	err error
}

// CallWithTimeout เรียก fn แต่รอไม่เกิน d
// ถ้าเกินเวลา ให้คืน error ที่ errors.Is(err, context.DeadlineExceeded) == true
func CallWithTimeout(ctx context.Context, d time.Duration, fn func(context.Context) (string, error)) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, d)
	defer cancel() // คืน resource ของ timer เสมอ

	// buffer 1: ถ้าเรา timeout ไปแล้วไม่มีใครรออ่าน goroutine ก็ยังส่งได้แล้วจบ → ไม่ leak
	ch := make(chan result, 1)
	go func() {
		v, err := fn(ctx)
		ch <- result{v, err}
	}()

	select {
	case r := <-ch:
		return r.val, r.err
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// FirstSuccess รันทุก fn พร้อมกัน คืนผลของตัวแรกที่สำเร็จ แล้ว cancel ตัวที่เหลือ
// ถ้าล้มเหลวทั้งหมด คืน error ที่รวมทุก error (errors.Join)
func FirstSuccess(ctx context.Context, fns ...func(context.Context) (string, error)) (string, error) {
	if len(fns) == 0 {
		return "", ErrNoFuncs
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel() // เมื่อได้ผลแล้ว return → cancel ตัวที่ยังทำงานอยู่

	ch := make(chan result, len(fns)) // buffer เท่าจำนวนงาน → ทุก goroutine ส่งได้ไม่ค้าง
	for _, fn := range fns {
		go func() {
			v, err := fn(ctx)
			ch <- result{v, err}
		}()
	}

	var errs []error
	for range fns {
		r := <-ch
		if r.err == nil {
			return r.val, nil
		}
		errs = append(errs, r.err)
	}
	return "", errors.Join(errs...)
}
