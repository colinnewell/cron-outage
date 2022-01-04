package cron

import (
	"fmt"
	"strconv"
	"strings"
	"time"
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
	Comment    string
	User       string
}

func (l Line) InWindow(start time.Time, end time.Time) bool {
	found := checkRange(int(start.Month()), int(end.Month()), l.Month, 12, 1)
	if !found {
		return false
	}

	// FIXME: take a closer look at the boundaries on this.
	found = checkRange(start.Day(), end.Day(), l.DayOfMonth, 31, 1)
	if !found {
		return false
	}

	found = checkRange(int(start.Weekday()), int(end.Weekday()), l.DayOfWeek, 7, 0)
	if !found {
		return false
	}

	found = checkRange(start.Hour(), end.Hour(), l.Hour, 24, 0)
	if !found {
		return false
	}

	if end.Sub(start).Hours() >= 1 {
		// NOTE: could do this trick for greater than a day etc. too.
		// assume all minutes are good since the duration is > 1 hour
		return true
	}

	found = checkRange(start.Minute(), end.Minute(), l.Minute, 60, 0)
	if !found {
		return false
	}

	return true
}

func checkRange(start, end int, valid []int, max int, lowest int) bool {
	found := false
	for m := start; m <= end && !found; m = (m + 1) % max {
		if m < lowest {
			m = lowest
		}
		for _, om := range valid {
			if om == int(m) {
				found = true
				break
			}
		}
	}
	return found
}

func (l Line) Lines() []string {
	lines := []string{
		generateLine("minute", l.Minute),
		generateLine("hour", l.Hour),
		generateLine("day of month", l.DayOfMonth),
		generateLine("month", l.Month),
		generateLine("day of week", l.DayOfWeek),
		fmt.Sprintf("%-14s%s", "user", l.User),
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
	done := false
	for a := 0; a < len(line) && !done; a++ {
		switch line[a] {
		case '#':
			if start > 0 {
				l.Command = line[start:a]
			}
			l.Comment = line[a:]
			return &l
		case ' ':
			if start >= 0 {
				// flush this
				if i >= len(parts) {
					l.User = line[start:a]
					start = a
					for ; start < len(line); start++ {
						if line[start] != ' ' {
							break
						}
					}
					done = true
					break
				}
				*parts[i].decoded = decode(line[start:a], parts[i].min, parts[i].max)
				start = -1
				i++
			}
		default:
			if start == -1 {
				start = a
			}
			if i > len(parts) {
				done = true
				break
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
