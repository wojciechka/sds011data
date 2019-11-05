package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/wojciechka/pm25data/pkg/sds011"
	"github.com/wojciechka/pm25data/pkg/rawdatastorage"	
)

var (
	dataDirectory = flag.String("dataDirectory", ".", "directory to store data in")
	device = flag.String("device", "/dev/ttyUSB0", "device to use for reading data")
	test = flag.Bool("test", false, "enable test mode that uses fake sensor data")
)

func errorOut(err error) {
	fmt.Fprintf(os.Stderr, "error when running application: %+v", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	var reader sds011.SensorReader

	// initialize storage
	storage, err := rawdatastorage.NewRawdataStorage(*dataDirectory)
	if err != nil {
		errorOut(err)
	}

	// initialize sensor reader
	if *test {
		reader = sds011.NewTestSensorReader()
	} else {
		reader, err = sds011.NewSensorReader(*device)
		if err != nil {
			errorOut(err)
		}
	}

	handler := sds011.NewSensorDataHandler(reader, storage)

	err = handler.Loop()
	if err != nil {
		errorOut(err)
	}
}
