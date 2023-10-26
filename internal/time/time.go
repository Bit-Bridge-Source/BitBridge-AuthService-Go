package time

import "time"

// TimeSource is an interface representing a source to get the current time.
type TimeSource interface {
	Now() time.Time
}

type SystemTime struct{}

func NewSystemTime() *SystemTime {
	return &SystemTime{}
}

// Now implements TimeSource.
func (*SystemTime) Now() time.Time {
	return time.Now()
}

// Ensure SystemTime implements TimeSource.
var _ TimeSource = (*SystemTime)(nil)
