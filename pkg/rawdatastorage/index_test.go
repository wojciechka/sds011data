package rawdatastorage

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestEasliestDataPoint(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	w, err := NewRawdataStorage(tempdir)
	if err != nil {
		log.Fatal(err)
	}

	when, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-11-01 13:00:00", time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}

	laterWhen := when.Add(36 * time.Hour)

	testData := []byte{1, 2, 3, 4}
	w.WriteData(when, testData)
	w.WriteData(laterWhen, testData)

	earliestDataPoint, err := w.EarliestDataPoint()
	if err != nil {
		log.Fatalf("unable to determine earliest data point: %v", err)
	}

	if earliestDataPoint != when {
		t.Errorf("earliest data point is %v, not %v", earliestDataPoint, when)
	}
}

func TestNextDataPoint(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	w, err := NewRawdataStorage(tempdir)
	if err != nil {
		log.Fatal(err)
	}

	when, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-11-01 13:00:00", time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}

	laterWhen := when.Add(36 * time.Hour)

	testData := []byte{1, 2, 3, 4}
	w.WriteData(when, testData)
	w.WriteData(laterWhen, testData)

	nextDataPoint, err := w.NextDataPoint(when)
	if err != nil {
		log.Fatalf("unable to determine next data point: %v", err)
	}

	if nextDataPoint != laterWhen {
		t.Errorf("next data point is %v, not %v", nextDataPoint, laterWhen)
	}
}

func TestLatestDataPoint(t *testing.T) {
	tempdir, err := ioutil.TempDir("", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tempdir)

	w, err := NewRawdataStorage(tempdir)
	if err != nil {
		log.Fatal(err)
	}

	when, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-11-01 13:00:00", time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}

	laterWhen := when.Add(36 * time.Hour)

	testData := []byte{1, 2, 3, 4}
	w.WriteData(when, testData)
	w.WriteData(laterWhen, testData)

	latestDataPoint, err := w.LatestDataPoint()
	if err != nil {
		log.Fatalf("unable to determine last data point: %v", err)
	}

	if latestDataPoint != laterWhen {
		t.Errorf("latest data point is %v, not %v", latestDataPoint, laterWhen)
	}
}
