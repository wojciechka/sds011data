package rawdatastorage

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"
)

func TestReadData(t *testing.T) {
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
	startOfDay, err := time.ParseInLocation("2006-01-02 15:04:05", "2019-11-01 00:00:00", time.Now().Location())
	if err != nil {
		log.Fatal(err)
	}

	testData := []byte{1, 2, 3, 4}
	w.WriteData(when, testData)

	// read bytes from beginning of day until now
	emptyBytes := make([]byte, 13*3600*4)
	numRead, err := w.ReadData(startOfDay, emptyBytes)
	if err != nil {
		log.Fatal(err)
	}

	if numRead != len(emptyBytes) {
		t.Errorf("Unable to read %v bytes - %v read", len(emptyBytes), numRead)
	}

	for i := 0; i < len(emptyBytes); i += 4 {
		if bytes.Compare(emptyBytes[i:i+4], w.defaultChunk) != 0 {
			t.Errorf("Bytes at %v differ", i)
		}
	}

	// read the written bytes
	testDataBuf := make([]byte, len(testData))
	numRead, err = w.ReadData(when, testDataBuf)
	if err != nil {
		log.Fatal(err)
	}

	if numRead != len(testData) {
		t.Errorf("Unable to read %v bytes - %v read", len(testData), numRead)
	}

	if bytes.Compare(testData, testDataBuf) != 0 {
		t.Error("Test data bytes differ")
	}
}
