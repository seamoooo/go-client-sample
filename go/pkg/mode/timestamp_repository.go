package mode

import (
	"context"
)

// TimestampRepository is the interface of repository for timestamp.
type TimestampRepository interface {
	BatchGetTimestamp(ctx context.Context, input BatchGetTimestampInput) (BatchGetTimestampOutput, error)
}

// BatchGetTimestampInput is input of BatchGetTimestamp
type BatchGetTimestampInput struct {
	TimeRange TimeRange
}

// BatchGetTimestampOutput is output of BatchGetTimestamp
type BatchGetTimestampOutput struct {
	Timestamp string
}
