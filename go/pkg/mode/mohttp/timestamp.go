package mohttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/seamoooo/take-home-assignment/go/pkg/mode"
)

// TimestampRepository is a Repository for parkings.
type TimestampRepository struct {
	host   string
	client *http.Client
}

// NewTimestampRepository allocates and returns TimestampRepository.
func NewTimestampRepository(host string) *TimestampRepository {

	// Temporary setting value. Change according to the situation.
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &TimestampRepository{host: host, client: client}
}

// BatchGetTimestamp is return all time stamp with data
// This data are returned as plain text. Each line represents one data point,
// starting with a timestamp
// eg)
// 2021-03-04T03:45:14Z 110.8634
// 2021-03-04T03:46:30Z 110.7046
// 2021-03-04T03:47:27Z 110.5467
func (r *TimestampRepository) BatchGetTimestamp(ctx context.Context, input mode.BatchGetTimestampInput) (mode.BatchGetTimestampOutput, error) {
	const op = "batch-get-time-stamp"

	url := fmt.Sprintf("%s/data?begin=%s&end=%s",
		r.host,
		input.TimeRange.Start.Format(time.RFC3339),
		input.TimeRange.End.Format(time.RFC3339),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return mode.BatchGetTimestampOutput{}, &mode.Error{
			Message: "fail to create client",
			Op:      op,
			Err:     err,
		}
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return mode.BatchGetTimestampOutput{}, &mode.Error{
			Message: "fail to request timestamp",
			Op:      op,
			Err:     err,
		}
	}

	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return mode.BatchGetTimestampOutput{}, &mode.Error{
			Message: "fail to read responce",
			Op:      op,
			Err:     err,
		}
	}
	return mode.BatchGetTimestampOutput{Timestamp: string(data)}, nil
}
