package sds011

import (
	"time"

	"github.com/juju/errors"
)

// DataHandler is an interface for logic handling reading of data
type DataHandler interface {
	Loop() error
}

type sensorDataHandler struct {
	reader SensorReader
	writer DataWriter
}

// NewSensorDataHandler creates a DataHandler for sensor reading
func NewSensorDataHandler(reader SensorReader, writer DataWriter) DataHandler {
	return &sensorDataHandler{reader: reader, writer: writer}
}

func (c *sensorDataHandler) Loop() error {
	for {
		when := time.Now()
		buf, err := c.reader.ReadData()
		if err != nil {
			return errors.Trace(err)
		}

		// TODO: calculate checksum

		if len(buf) == 10 {
			err = c.writer.WriteData(when, buf[2:6])
			if err != nil {
				return errors.Trace(err)
			}
		}
	}
}
