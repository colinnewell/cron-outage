package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/colinnewell/cron-outage/pkg/cron"
	"github.com/spf13/pflag"
)

func main() {
	var displayVersion bool
	var start, end, notBefore, timeFormat string
	// FIXME: don't sort the flags in the help
	pflag.StringVar(&start, "start", "", "Start of outage")
	pflag.StringVar(&end, "end", "", "End of outage")
	pflag.StringVar(&notBefore, "not-before", "", "Also check if it won't run before this time (optional)")
	pflag.StringVar(&timeFormat, "time-format", "2006-01-02T15:04:05Z", "Time format to parse")
	pflag.BoolVar(&displayVersion, "version", false, "Display program version")
	pflag.Parse()

	if displayVersion {
		fmt.Println("Version:", Version)
		return
	}

	if timeFormat == "" || start == "" || end == "" {
		pflag.Usage()
		return
	}

	s, err := time.Parse(timeFormat, start)
	if err != nil {
		fmt.Printf("Error parsing %s: %s\n", s, err)
		pflag.Usage()
		return
	}

	e, err := time.Parse(timeFormat, end)
	if err != nil {
		fmt.Printf("Error parsing %s: %s\n", e, err)
		pflag.Usage()
		return
	}

	var nb time.Time
	if notBefore != "" {
		var err error
		nb, err = time.Parse(timeFormat, notBefore)
		if err != nil {
			fmt.Printf("Error parsing %s: %s\n", e, err)
			pflag.Usage()
			return
		}
	}

	for _, v := range pflag.Args() {
		// we aren't trying to constrain the paths so we make use
		// of what the user gives us as is.
		f, err := os.Open(v)
		if err != nil {
			fmt.Printf("Unable to open %s: %s\n", v, err)
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := scanner.Text()
			lines := cron.ParseLine(line)
			if lines != nil && lines.Command != "" && lines.InWindow(s, e) {
				if nb.IsZero() || !lines.InWindow(e, nb) {
					fmt.Println(line)
				}
			}
		}
	}
}
