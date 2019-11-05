package sds011

import (
	"time"
)

// DataReader defines an interface for reading data from storage
type DataReader interface {
	ReadData(when time.Time, buf []byte) (int, error)
	EarliestDataPoint() (time.Time, error)
	NextDataPoint(after time.Time) (time.Time, error)
	LatestDataPoint() (time.Time, error)
}

// DataWriter defines an interface for writing data to storage
type DataWriter interface {
	WriteData(when time.Time, buf []byte) error
}
