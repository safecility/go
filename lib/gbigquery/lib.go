package gbigquery

import (
	"fmt"
	"time"
)

type TimeInterval string

const (
	HOUR TimeInterval = "HOUR"
	DAY  TimeInterval = "DAY"
)

type QueryInterval struct {
	Start time.Time `json:"start,omitempty"`
	End   time.Time `json:"end,omitempty"`
}

type BucketType struct {
	Interval   TimeInterval
	Multiplier int
}

func (bt BucketType) String() string {
	return fmt.Sprintf("%v-%v(s)", bt.Multiplier, bt.Interval)
}

type Bucket struct {
	StartTime time.Time     `datastore:"Time"`
	Duration  time.Duration `datastore:"-"`
}
