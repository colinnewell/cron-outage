package cron_test

import (
	"testing"
	"time"

	"github.com/colinnewell/cron-outage/internal/cron"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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
			"user          ",
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

func TestInWindow(t *testing.T) {
	l := cron.Line{
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

	s := time.Date(2020, 11, 17, 10, 34, 58, 0, time.UTC)
	e := time.Date(2020, 11, 17, 11, 34, 58, 0, time.UTC)
	assert.False(t, l.InWindow(s, e), "Should not be in window")

	s = time.Date(2020, 11, 17, 4, 3, 58, 0, time.UTC)
	e = time.Date(2020, 11, 17, 4, 5, 58, 0, time.UTC)
	assert.True(t, l.InWindow(s, e), "Should be in window")
}
