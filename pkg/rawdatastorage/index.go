package rawdatastorage

import (
	"bytes"
	"io"
	"io/ioutil"
	"regexp"
	"sort"
	"time"

	"github.com/juju/errors"
)

var (
	filenamePattern    = regexp.MustCompile("^[0-9]{8}$")
	filenameDateFormat = "20060102"

	// ErrUnableToFindData indicates an error about not being able to find data or further data in a directory
	ErrUnableToFindData = errors.New("unable to find any data")
)

// EarliestDataPoint finds first day and first timestamp on that day that has any data available or returns an error if no data is available at all
func (w *RawdataStorage) EarliestDataPoint() (time.Time, error) {
	return w.NextDataPoint(time.Time{})
}

// NextDataPoint finds next day and first timestamp on that day that has any data available or returns an error if no data is available at all
func (w *RawdataStorage) NextDataPoint(after time.Time) (time.Time, error) {
	sinceStr := after.AddDate(0, 0, 1).Format(filenameDateFormat)

	files, err := w.readDir(false, func(name string) bool {
		return name >= sinceStr
	})
	if err != nil {
		return time.Time{}, errors.Trace(err)
	}

	return w.findOffsetInFiles(files, firstOffset)
}

// LatestDataPoint finds last day and latest timestamp on that day that has any data available or returns an error if no data is available at all
func (w *RawdataStorage) LatestDataPoint() (time.Time, error) {
	files, err := w.readDir(true, func(name string) bool {
		return true
	})
	if err != nil {
		return time.Time{}, errors.Trace(err)
	}

	return w.findOffsetInFiles(files, firstOffset)
}

func filenameToTime(name string) (time.Time, error) {
	return time.ParseInLocation(filenameDateFormat, name, time.Now().Location())
}

func (w *RawdataStorage) readDir(reverse bool, filter func(name string) bool) ([]string, error) {
	files, err := ioutil.ReadDir(w.dataDirectory)
	if err != nil {
		return []string{}, errors.Trace(err)
	}

	var res []string

	for _, fileInfo := range files {
		name := fileInfo.Name()
		if len(name) == 8 && filenamePattern.MatchString(name) {
			if filter(name) {
				res = append(res, name)
			}
		}
	}

	if len(res) == 0 {
		return []string{}, ErrUnableToFindData
	}

	sort.Strings(res)

	if reverse {
		// reverse the array
		for i1 := 0; i1 < len(res)/2; i1++ {
			i2 := len(res) - i1 - 1
			res[i1], res[i2] = res[i2], res[i1]
		}
	}

	return res, nil
}

type findOffsetOption int

const (
	firstOffset findOffsetOption = 1
	lastOffset  findOffsetOption = 2
)

func (w *RawdataStorage) findOffsetInFiles(files []string, what findOffsetOption) (time.Time, error) {
	for _, filename := range files {
		when, err := filenameToTime(filename)
		if err != nil {
			return time.Time{}, errors.Trace(err)
		}

		until := when.AddDate(0, 0, 1)
		when, found, err := w.findOffset(when, until, firstOffset)

		if found {
			return when, nil
		}
	}

	return time.Time{}, ErrUnableToFindData
}

func (w *RawdataStorage) findOffset(when time.Time, until time.Time, what findOffsetOption) (time.Time, bool, error) {
	found := false

	for now := when; now.Before(until); now = now.Add(1 * time.Second) {
		var buf [4]byte
		numBytes, err := w.ReadData(now, buf[:])
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return time.Time{}, false, errors.Trace(err)
			}
		}
		if numBytes == 0 {
			break
		} else if numBytes != 4 {
			return time.Time{}, false, errors.New("unable to fully read data")
		}

		if bytes.Compare(buf[:], w.defaultChunk) != 0 {
			found = true
			when = now
			if what == firstOffset {
				break
			}
		}
	}

	return when, found, nil
}
