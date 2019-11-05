package main

import (
	"os"
	"strings"
	"time"

	"github.com/juju/errors"
	"github.com/wojciechka/pm25data/pkg/rawdatastorage"
)

func getResumePoint(stateFile string, storage *rawdatastorage.RawdataStorage) (time.Time, error) {
	file, err := os.OpenFile(stateFile, os.O_RDONLY, 0o777)
	if err != nil {
		if os.IsNotExist(err) {
			return storage.EarliestDataPoint()
		}
		return time.Time{}, errors.Trace(err)
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return time.Time{}, errors.Trace(err)
	}

	if fileInfo.Size() > 64 {
		return time.Time{}, errors.New("unable to read state file")
	}

	buffer := make([]byte, fileInfo.Size())
	byteCount, err := file.Read(buffer)
	if err != nil {
		return time.Time{}, errors.Trace(err)
	}

	if int64(byteCount) != fileInfo.Size() {
		return time.Time{}, errors.New("unable to read state file")
	}

	return time.Parse(time.RFC3339, strings.TrimSpace(string(buffer)))
}

func writeResumePoint(stateFile string, when time.Time) error {
	file, err := os.Create(stateFile)
	if err != nil {
		return errors.Trace(err)
	}

	defer file.Close()

	_, err = file.WriteString(when.Format(time.RFC3339))
	return err
}
