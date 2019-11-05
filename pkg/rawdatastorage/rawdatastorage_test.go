package rawdatastorage

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"
)

func TestSeekToTime(t *testing.T) {
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

	w.seekToTime(when, true)

	info, err := os.Stat(path.Join(tempdir, when.Format("20060102")))
	if err != nil {
		t.Errorf("File %s not found or could not be checked: %v", when.Format("20060102"), err)
	}

	if info.Size() != (13 * 3600 * 4) {
		t.Errorf("Invalid size for timestamp file: %v", info.Size())
	}
}
