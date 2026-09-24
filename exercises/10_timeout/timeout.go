// Package timeout — ข้อ 10: context.WithTimeout, select, goroutine leak
package timeout

import (
	"context"
	"errors"
	"time"
)

var ErrNoFuncs = errors.New("no functions given")

// CallWithTimeout เรียก fn แต่รอไม่เกิน d
// ถ้าเกินเวลา ให้คืน error ที่ errors.Is(err, context.DeadlineExceeded) == true
// และต้องไม่มี goroutine ค้าง แม้ fn จะไม่สนใจ ctx ก็ตาม
func CallWithTimeout(ctx context.Context, d time.Duration, fn func(context.Context) (string, error)) (string, error) {
	// TODO: ตอนนี้รอ fn จนเสร็จโดยไม่สน timeout
	return fn(ctx)
}

// FirstSuccess รันทุก fn พร้อมกัน คืนผลของตัวแรกที่สำเร็จ แล้ว cancel ตัวที่เหลือ
// ถ้าล้มเหลวทั้งหมด คืน error ที่รวมทุก error, ไม่มี fn เลย → ErrNoFuncs
func FirstSuccess(ctx context.Context, fns ...func(context.Context) (string, error)) (string, error) {
	// TODO
	return "", nil
}
