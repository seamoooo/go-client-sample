package mode

import (
	"bufio"
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SummarizeHourlyTimestamp summarizing hourly time stamp
type SummarizeHourlyTimestamp struct {
	TimestampRepository TimestampRepository
}

// SummarizeHourlyTimestampInput is input of Do method
type SummarizeHourlyTimestampInput struct {
	TimeRange TimeRange
}

// SummarizeHourlyTimestampOutput is output of Do method
type SummarizeHourlyTimestampOutput struct {
	SummarizeHourlyTimestamp string
}

// Do is summarizing hourly time stamp
func (t *SummarizeHourlyTimestamp) Do(
	ctx context.Context,
	in SummarizeHourlyTimestampInput,
) (SummarizeHourlyTimestampOutput, error) {
	const op = "summarize-hourly-time-stamp"

	out, err := t.TimestampRepository.BatchGetTimestamp(ctx, BatchGetTimestampInput(in))
	if err != nil {
		return SummarizeHourlyTimestampOutput{}, &Error{
			Op:      op,
			Message: "fail to batch get timestamp",
			Err:     err,
		}
	}

	hourlyData, err := calculateHourlyAverages(out.Timestamp, in.TimeRange)
	if err != nil {
		return SummarizeHourlyTimestampOutput{}, &Error{
			Op:      op,
			Message: "fail to calculate hourly averages",
			Err:     err,
		}
	}

	serializeData := serializeHourlyDataSorted(hourlyData)
	return SummarizeHourlyTimestampOutput{SummarizeHourlyTimestamp: serializeData}, nil
}

func calculateHourlyAverages(data string, timeRange TimeRange) (map[time.Time]float64, error) {
	const op = "calculate-hourly-averages"

	hourlyData := make(map[time.Time]float64)
	var currentBucketTs time.Time
	var currentBucketValues []float64

	assignDataPointToHourlyBucket := func(ts time.Time, value float64) {
		bucketTs := time.Date(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), 0, 0, 0, ts.Location())
		if bucketTs != currentBucketTs {
			if len(currentBucketValues) > 0 {
				hourlyData = updateHourlyData(currentBucketTs, currentBucketValues, hourlyData)
			}
			currentBucketTs = bucketTs
			currentBucketValues = make([]float64, 0)
		}
		currentBucketValues = append(currentBucketValues, value)
	}

	scanner := bufio.NewScanner(strings.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		ts, value, err := parseTimestampLine(line)
		if err != nil {
			return nil, err
		}

		if !timeRange.Contains(ts.UTC()) {
			return nil, &Error{
				Op:      op,
				Err:     err,
				Message: fmt.Sprintf("timestamp is out of range: %v", ts.UTC()),
			}
		}

		assignDataPointToHourlyBucket(ts, value)
	}

	if len(currentBucketValues) > 0 {
		hourlyData = updateHourlyData(currentBucketTs, currentBucketValues, hourlyData)
	}

	return hourlyData, nil
}

func parseTimestampLine(line string) (time.Time, float64, error) {
	const op = "parse-timestamp-line"

	// This regular expression matches strings that follow this pattern
	// eg) 2021-03-04T05:24:05Z 100.4288
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)\s+(\S+)$`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 3 {
		return time.Time{}, 0, &Error{
			Op:      op,
			Message: fmt.Sprintf("invalid timestamp line: %s", line),
		}
	}

	ts, err := time.Parse(time.RFC3339, matches[1])
	if err != nil {
		return time.Time{}, 0,
			&Error{
				Op:      op,
				Message: fmt.Sprintf("fail to parse string to time: %s", matches[1]),
			}
	}

	value, err := strconv.ParseFloat(matches[2], 64)
	if err != nil {
		return time.Time{}, 0,
			&Error{
				Op:      op,
				Message: fmt.Sprintf("fail to parse value string to float: %s", matches[1]),
			}
	}

	return ts, value, nil
}

func updateHourlyData(
	bucketTs time.Time,
	bucketValues []float64,
	hourlyData map[time.Time]float64,
) map[time.Time]float64 {
	sum := 0.0
	for _, v := range bucketValues {
		sum += v
	}
	currentAvg, ok := hourlyData[bucketTs.UTC()]
	if !ok {
		hourlyData[bucketTs.UTC()] = sum / float64(len(bucketValues))
	} else {
		newAvg := (currentAvg*float64(len(bucketValues)) + sum) / float64(len(bucketValues)*2)
		hourlyData[bucketTs.UTC()] = newAvg
	}
	return hourlyData
}

func serializeHourlyDataSorted(m map[time.Time]float64) string {
	var sortedKeys []time.Time
	for k := range m {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Slice(sortedKeys, func(i, j int) bool {
		return sortedKeys[i].Before(sortedKeys[j])
	})

	var sb strings.Builder
	for _, k := range sortedKeys {
		sb.WriteString(fmt.Sprintf("%s %.4f\n", k.Format(time.RFC3339), m[k]))
	}
	return sb.String()
}
