package mode

import "time"

// TimeRange represents a range of time (Start End).
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// Empty returns the time range is empty or not.
func (r TimeRange) Empty() bool {
	return r == TimeRange{}
}

// Contains returns true iff t is in (Start, End).
func (r TimeRange) Contains(t time.Time) bool {
	return (t.Equal(r.Start) || t.After(r.Start)) && (t.Equal(r.End) || t.Before(r.End))
}

// ParseTimeRange parses the command line arguments and returns a TimeRange.
func ParseTimeRange(startTime string, endTime string) (TimeRange, error) {
	const op = "parse-time-range"

	begin, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return TimeRange{}, &Error{
			Code:    ErrInvalidArgument,
			Message: "Invalid begin timestamp",
			Op:      op,
			Err:     err,
		}
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return TimeRange{}, &Error{
			Code:    ErrInvalidArgument,
			Message: "Invalid end timestamp",
			Op:      op,
			Err:     err,
		}
	}

	return TimeRange{Start: begin, End: end}, nil
}
