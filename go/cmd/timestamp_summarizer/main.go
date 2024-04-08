package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/seamoooo/take-home-assignment/go/pkg/mode"
	"github.com/seamoooo/take-home-assignment/go/pkg/mode/mohttp"
)

const modeAPIhost = "https://tsserv.tinkermode.dev"

func main() {
	if len(os.Args) != 3 {
		fatalError("Usage: %s <begin_timestamp> <end_timestamp>", os.Args[0])
	}

	timeRange, err := mode.ParseTimeRange(os.Args[1], os.Args[2])
	if err != nil {
		fatalError("%v", err)
	}

	repo := mohttp.NewTimestampRepository(modeAPIhost)
	input := mode.SummarizeHourlyTimestampInput{TimeRange: timeRange}
	sht := mode.SummarizeHourlyTimestamp{
		TimestampRepository: repo,
	}

	ctx := context.Background()

	out, err := sht.Do(ctx, input)
	if err != nil {
		fatalError("%v", err)
	}

	fmt.Printf("%s", out.SummarizeHourlyTimestamp)
}

func fatalError(format string, args ...interface{}) {
	log.Printf(format, args...)
	os.Exit(1)
}
