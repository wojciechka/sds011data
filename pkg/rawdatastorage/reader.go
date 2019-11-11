package rawdatastorage

import (
	"io"
	"time"

	"github.com/juju/errors"
)

// ReadData reads data available at specified
func (w *RawdataStorage) ReadData(when time.Time, bytes []byte) (int, error) {
	offset, err := w.seekToTime(when, false)
	if err != nil {
		w.closeCurrentFile()
		// return EOF so it can be handled properly
		if err == io.EOF {
			return 0, err
		}
		return 0, errors.Trace(err)
	}

	numRead, err := w.currentFile.Read(bytes)
	if err != nil {
		w.closeCurrentFile()
		// return EOF so it can be handled properly
		if err == io.EOF {
			return 0, err
		}
		return 0, errors.Trace(err)
	}

	w.currentOffset = offset + int64(numRead)
	return numRead, nil
}
