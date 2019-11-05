package rawdatastorage

import (
	"time"

	"github.com/juju/errors"
)

// WriteData writes data for specified date and time
func (w *RawdataStorage) WriteData(when time.Time, buf []byte) error {
	offset, err := w.seekToTime(when, true)
	if err != nil {
		w.closeCurrentFile()
		return errors.Trace(err)
	}

	_, err = w.currentFile.Write(buf)
	if err != nil {
		w.closeCurrentFile()
		return errors.Trace(err)
	}

	w.currentOffset = offset + int64(len(buf))

	return nil
}
