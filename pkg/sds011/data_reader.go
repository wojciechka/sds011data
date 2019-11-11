package sds011

import (
	"encoding/binary"
	"io"
	"time"

	"github.com/juju/errors"
	"github.com/wojciechka/pm25data/pkg/rawdatastorage"
)

// DataValues defines PM2.5 and PM10 raw or average values for Sds011 output
type DataValues struct {
	When  time.Time
	Count int
	Pm25  int
	Pm10  int
}

// DataParser defines interface for reading data, possibly also grouping it
type DataParser interface {
	ReadData() (*DataValues, error)
	Position() time.Time
}

type sds011GroupDataParser struct {
	reader     DataReader
	position   time.Time
	groupBy    time.Duration
	seconds    int
	buffer     []byte
	pm25buffer []int16
	pm10buffer []int16
}

// NewDataParser creates a new DataParser based on parameters
func NewDataParser(reader DataReader, groupBy time.Duration, startAt time.Time) DataParser {
	seconds := int(groupBy / time.Second)
	return &sds011GroupDataParser{
		reader:   reader,
		position: startAt.Truncate(groupBy),
		groupBy:  groupBy,
		seconds:  seconds,
		buffer:   make([]byte, seconds*4),
	}
}

func (p *sds011GroupDataParser) ReadData() (*DataValues, error) {
	result := &DataValues{When: p.position, Count: 0, Pm25: 0, Pm10: 0}

	numBytes, err := p.reader.ReadData(p.position, p.buffer)
	if err != nil {
		if err == io.EOF {
			// if EOF was received received, check if there is any data after today
			newPosition, err := p.reader.NextDataPoint(p.position)
			if err != nil {
				if err == rawdatastorage.ErrUnableToFindData {
					// to nothing, hopefully more data will be available soon
					return result, nil
				}
			}
			p.position = newPosition
			return result, nil
		}
		return result, errors.Trace(err)
	}

	numSeconds := numBytes / 4
	if numSeconds == 0 {
		return result, nil
	}

	for i := 0; i < numSeconds; i++ {
		pm25 := binary.LittleEndian.Uint16(p.buffer[i*4+0 : i*4+2])
		pm10 := binary.LittleEndian.Uint16(p.buffer[i*4+2 : i*4+4])

		if pm10 != 0xffff && pm25 != 0xffff {
			result.Count++
			result.Pm25 += int(pm25)
			result.Pm10 += int(pm10)
		}

	}

	p.position = p.position.Add(time.Duration(numSeconds) * time.Second)
	if result.Count > 0 {
		result.Pm25 = result.Pm25 / result.Count
		result.Pm10 = result.Pm10 / result.Count
	}
	return result, nil
}

func (p *sds011GroupDataParser) Position() time.Time {
	return p.position
}
