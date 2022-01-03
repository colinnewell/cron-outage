package main

import (
	"flag"
	"fmt"

	"github.com/colinnewell/cron-outage/internal/cron"
)

func main() {
	var displayVersion bool
	flag.BoolVar(&displayVersion, "version", false, "Display program version")
	flag.Parse()

	if displayVersion {
		fmt.Println("Version:", Version)
		return
	}

	for _, v := range flag.Args() {
		lines := cron.ParseLine(v)
		for _, l := range lines.Lines() {
			fmt.Println(l)
		}
	}
}
