package snowflake

import (
	"time"
)

// Overflow struct holds the details about how much a Node has overflown while
// generating IDs.
type Overflow struct {
	// Time at which the overflow was read from Node
	Time time.Time
	// Ms is the milliseconds by which the Node is ahead of now.
	// Ms in combination with Step gives us exactly how many IDs
	// have overflown.
	Ms int64
	// Step is the number of steps by which the Node is ahead of now.
	// Step in combination with Ms gives us exactly how many IDs
	// have overflown.
	Step int64
}

// IsZero evaluates to true if there is no overflow
func (o Overflow) IsZero() bool {
	return o.Ms == 0
}

// Duration gives the time-duration equivalent of the overflow.
// The returned value can be used to wait till the overflow is "gone". Note
// that waiting will not always work if more than one goroutine are causing
// continuous overflow.
func (o Overflow) Duration() time.Duration {
	if o.IsZero() {
		return 0
	}

	return time.Duration(o.Ms) * time.Millisecond
}

func (o Overflow) DurationToClear() time.Duration {
	d := o.Duration()
	if d <= 0 {
		return 0
	}

	now := time.Now()
	if o.Time.After(now) {
		return 0
	}

	return d - now.Sub(o.Time)
}

// AfterCleared can be used to suspend the current goroutine until the overflow
// is cleared. Note that the overflown Node might still be overflowing if
// multiple goroutines are generating ids from the same Node.
func (o Overflow) AfterCleared() <-chan time.Time {
	return time.After(o.DurationToClear())
}
