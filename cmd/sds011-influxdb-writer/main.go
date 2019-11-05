package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/influxdata/influxdb-client-go"

	"github.com/wojciechka/pm25data/pkg/rawdatastorage"
	"github.com/wojciechka/pm25data/pkg/sds011"
)

var (
	dataDirectory = flag.String("dataDirectory", ".", "directory to store data in")
	stateFile     = flag.String("stateFile", "./.upload-to-influx.state", "file to keep state of last uploaded data")
	influxURL     = flag.String("influxUrl", "https://us-west-2-1.aws.cloud2.influxdata.com/", "URL address of InfluxDB endpoint")
	influxOrg     = flag.String("influxOrg", "", "organization to send data to")
	influxBucket  = flag.String("influxBucket", "", "bucket to send data to")
	influxName    = flag.String("influxName", "sds011", "name of the metric to upload as")
	influxTags    = flag.String("influxTags", "", "comma separated tags to apply")
	influxToken   = flag.String("influxToken", "", "token to use for authentication to InfluxDB")
	groupBy       = flag.Duration("groupBy", 1*time.Second, "amount of time by which to group the data by")
	wait          = flag.Bool("wait", false, "wait for data to become available")
)

func errorOut(err error) {
	fmt.Fprintf(os.Stderr, "error when running application: %+v", err)
	os.Exit(1)
}

func main() {
	flag.Parse()

	// initialize storage
	storage, err := rawdatastorage.NewRawdataStorage(*dataDirectory)
	if err != nil {
		errorOut(err)
	}

	when, err := getResumePoint(*stateFile, storage)
	if err != nil {
		errorOut(err)
	}

	parser := sds011.NewDataParser(storage, *groupBy, when)

	until, err := storage.LatestDataPoint()
	if err != nil {
		errorOut(err)
	}

	influx, err := influxdb.New(*influxURL, *influxToken)
	if err != nil {
		errorOut(err)
	}
	defer influx.Close()

	influxMetrics := []influxdb.Metric{}

	writeMetrics := func() {
		if len(influxMetrics) > 0 {
			n, err := influx.Write(context.Background(), *influxBucket, *influxOrg, influxMetrics...)
			if err != nil {
				errorOut(err)
			}
			err = writeResumePoint(*stateFile, parser.Position())
			if err != nil {
				errorOut(err)
			}
			fmt.Printf("%d of %d metric(s) sent; currently at %s\n", n, len(influxMetrics), parser.Position().Format(time.RFC1123))
			influxMetrics = []influxdb.Metric{}
		}
	}

	influxTagMap := toTags(*influxTags)

	for true {
		d, err := parser.ReadData()
		if err != nil {
			if err == io.EOF {
				break
			}
			errorOut(err)
		}

		if d.Count > 0 {
			influxMetrics = append(influxMetrics, influxdb.NewRowMetric(
				map[string]interface{}{"pm25": float64(d.Pm25) / 10.0, "pm10": float64(d.Pm10) / 10.0},
				*influxName,
				influxTagMap,
				d.When,
			))

			// group up to 1000 metrics before sending them
			if len(influxMetrics) >= 1000 {
				writeMetrics()
			}
		} else {
			writeMetrics()

			// if the data was already past our data, let's re-read current latest state
			if d.When.After(until) {
				until, err = storage.LatestDataPoint()
				if err != nil {
					errorOut(err)
				}
			}

			// if still went over our available data, exit or wait
			if d.When.After(until) {
				if *wait {
					time.Sleep(*groupBy)
				} else {
					break
				}
			}
		}
	}

	writeMetrics()
}
