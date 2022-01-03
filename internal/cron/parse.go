package cron

import (
	"fmt"
	"strconv"
	"strings"
)

type Line struct {
	// FIXME: could add attributes that specify range and heading names
	// possibly point out validation errors?
	Minute     []int
	Hour       []int
	DayOfMonth []int
	Month      []int
	DayOfWeek  []int
	Command    string
}

func (l Line) Lines() []string {
	lines := []string{
		generateLine("minute", l.Minute),
		generateLine("hour", l.Hour),
		generateLine("day of month", l.DayOfMonth),
		generateLine("month", l.Month),
		generateLine("day of week", l.DayOfWeek),
		fmt.Sprintf("%-14s%s", "command", l.Command),
	}
	return lines
}

func generateLine(heading string, numbers []int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-14s", heading))
	for i, v := range numbers {
		if i > 0 {
			sb.WriteString(" ")
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}

func ParseLine(line string) *Line {
	// grab line info.
	// allow comments
	l := Line{}
	parts := []struct {
		min     int
		max     int
		decoded *[]int
	}{
		{0, 59, &l.Minute},
		{0, 23, &l.Hour},
		{1, 31, &l.DayOfMonth},
		{1, 12, &l.Month},
		{0, 6, &l.DayOfWeek},
	}
	i := 0
	start := -1
	for a := 0; a < len(line); a++ {
		switch line[a] {
		case ' ':
			if start >= 0 {
				// flush this
				*parts[i].decoded = decode(line[start:a], parts[i].min, parts[i].max)
				start = -1
				i++
			}
		default:
			if start == -1 {
				start = a
				if i >= len(parts) {
					break
				}
			}
		}
	}
	if start > 0 {
		l.Command = line[start:]
	}
	// */15
	// 0,15
	// 1-5
	// 0
	// *
	return &l
}

// decode takes the encoded values and turns them into an
// array of numbers
// light on error handling.  Just silently tries to do the
// best it can.
func decode(text string, min, max int) []int {
	tokens := tokenise(text)
	p := parser{min: min, max: max, tokens: tokens}
	p.processList()
	return p.vals
}