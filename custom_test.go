package deepequal_test

import (
	"testing"
	"time"

	"github.com/powerman/deepequal"
)

type Time time.Time

func (Time) Equal(time.Time) bool { return true } // Invalid signature.

func TestDeepEqualEqual(t *testing.T) {
	type T struct {
		t1 time.Time
		t2 *time.Time
	}
	var (
		zero     time.Time
		now      = time.Now().In(time.FixedZone("test", 3600))
		now2     = now.UTC()
		nowTime  = Time(now)
		now2Time = Time(now2)
	)

	tests := []struct {
		a, b interface{}
		want bool
	}{
		{now, now, true},
		{now, zero, false},
		{&now, now, false},
		{&now, &now, true},
		{now, now2, true},
		{&now, &now2, true},
		{T{now, &now}, T{now2, &now2}, true},
		{T{now, &now}, T{now2, &zero}, false},
		{nowTime, now, false},
		{nowTime, nowTime, true},
		{nowTime, now2Time, false},
	}
	for _, tc := range tests {
		tc := tc
		t.Run("", func(t *testing.T) {
			if res := deepequal.DeepEqual(tc.a, tc.b); res != tc.want {
				t.Errorf("DeepEqual(%v, %v) = %v, want %v", tc.a, tc.b, res, tc.want)
			}
		})
	}
}
