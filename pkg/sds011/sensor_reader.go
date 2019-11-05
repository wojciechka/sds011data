package sds011

import (
	"math/rand"
	"time"

	"github.com/juju/errors"
	"github.com/tarm/serial"
)

// SensorReader is an interface for reading data from a storage
type SensorReader interface {
	ReadData() ([]byte, error)
}

type realSensorReader struct {
	portConfig *serial.Config
	port       *serial.Port
}

// testSensorReader is an implementation of SensorReader that generates random data
type testSensorReader struct {
}

// NewSensorReader creates a new instance of SensorReader that reads data from serial port
func NewSensorReader(deviceName string) (SensorReader, error) {
	rc := &realSensorReader{}
	rc.portConfig = &serial.Config{
		Name: deviceName,
		Baud: 9600,
	}

	port, err := serial.OpenPort(rc.portConfig)
	if err != nil {
		return nil, errors.Trace(err)
	}

	rc.port = port
	return rc, nil
}

func (s *realSensorReader) ReadData() ([]byte, error) {
	buf := make([]byte, 10)
	n, err := s.port.Read(buf)
	if err != nil {
		return []byte{}, err
	}

	// TODO: handle the case where reading from sensor was not fully completed
	if n < 10 {
		return []byte{}, nil
	}
	return buf, nil
}

// NewTestSensorReader creates a new instance of SensorReader that generates random data
func NewTestSensorReader() SensorReader {
	return &testSensorReader{}
}

func (s *testSensorReader) ReadData() ([]byte, error) {
	time.Sleep(1 * time.Second)
	buf := []byte{
		0, 0,
		byte((rand.Int() & 0xff)),
		byte((rand.Int() & 0x0f)),
		byte((rand.Int() & 0xff)),
		byte((rand.Int() & 0x0f)),
		0, 0, 0, 0,
	}
	return buf, nil
}
