package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/colinnewell/cron-outage/internal/cron"
	"github.com/spf13/pflag"
)

func main() {
	var displayVersion bool
	var start, end, timeFormat string
	pflag.BoolVar(&displayVersion, "version", false, "Display program version")
	pflag.StringVar(&start, "start", "", "Start of outage")
	pflag.StringVar(&end, "end", "", "End of outage")
	pflag.StringVar(&timeFormat, "time-format", "2006-01-02T15:04:05Z", "Time format to parse")
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

	for _, v := range pflag.Args() {
		f, err := os.Open(v)
		if err != nil {
			fmt.Printf("Unable to open %s: %s\n", v, err)
			continue
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			lines := cron.ParseLine(scanner.Text())
			if lines != nil && lines.Command != "" {
				for _, l := range lines.Lines() {
					fmt.Println(l)
				}
			}
		}

	}
}
