package cron_test

import (
	"testing"

	"github.com/colinnewell/cron-outage/internal/cron"

	"github.com/google/go-cmp/cmp"
)

func TestOutput(t *testing.T) {
	l := cron.Line{
		Minute:  []int{3, 4},
		Hour:    []int{4},
		Command: "test",
	}

	if diff := cmp.Diff(
		l.Lines(),
		[]string{
			"minute        3 4",
			"hour          4",
			"day of month  ",
			"month         ",
			"day of week   ",
			"command       test",
		},
	); diff != "" {
		t.Errorf("Different output (-got +expected):\n%s\n", diff)
	}
}

func TestParse(t *testing.T) {
	expected := cron.Line{
		Minute: []int{3, 4},
		Hour:   []int{4},
		DayOfMonth: []int{
			1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14,
			15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26,
			27, 28, 29, 30, 31,
		},
		Month:     []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		DayOfWeek: []int{0, 1, 2, 3, 4, 5, 6},
		Command:   "test",
	}
	if diff := cmp.Diff(
		cron.ParseLine("3,4 4 * * *    test"), &expected,
	); diff != "" {
		t.Errorf("Different output (-got +expected):\n%s\n", diff)
		t.Logf("%#v", cron.ParseLine("3,4 4 * * *    test"))
	}
}
