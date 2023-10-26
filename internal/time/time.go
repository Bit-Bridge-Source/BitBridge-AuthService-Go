package time

import "time"

// TimeSource is an interface representing a source to get the current time.
type TimeSource interface {
	Now() time.Time
}
