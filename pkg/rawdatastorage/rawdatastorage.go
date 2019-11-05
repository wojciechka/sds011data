package rawdatastorage

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/juju/errors"
)

// RawdataStorage keeps sensor results as flat files on disk
type RawdataStorage struct {
	dataDirectory   string
	defaultChunk    []byte
	currentFilename string
	currentFile     *os.File
	currentOffset   int64
}

// NewRawdataStorage creates an instance of RawdataStorage storing data in a specified directory
func NewRawdataStorage(dataDirectory string) (*RawdataStorage, error) {
	defaultChunk := []byte{0xff, 0xff, 0xff, 0xff}

	// TODO: create directory if needed?

	return &RawdataStorage{
		dataDirectory: dataDirectory,
		defaultChunk:  defaultChunk,
	}, nil
}

func (w *RawdataStorage) ensureFileOpen(filename string) error {
	if w.currentFilename != filename {
		w.closeCurrentFile()
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return errors.Trace(err)
		}

		o, err := f.Seek(0, os.SEEK_END)
		if err != nil {
			return errors.Trace(err)
		}

		w.currentFile = f
		w.currentFilename = filename
		w.currentOffset = o
	}

	return nil
}

func (w *RawdataStorage) seekToTime(when time.Time, extend bool) (int64, error) {
	filename := path.Join(w.dataDirectory, when.Format(filenameDateFormat))
	offset := int64((when.Second() + when.Minute()*60 + when.Hour()*3600) * len(w.defaultChunk))

	err := w.ensureFileOpen(filename)
	if err != nil {
		return 0, errors.Trace(err)
	}

	if w.currentOffset != offset {
		// if true {
		o, err := w.currentFile.Seek(0, os.SEEK_END)
		if err != nil {
			w.closeCurrentFile()
			return 0, errors.Trace(err)
		}

		if o < offset {
			if extend {
				for times := (offset - o) / int64(len(w.defaultChunk)); times > 0; times-- {
					_, err := w.currentFile.Write(w.defaultChunk)
					if err != nil {
						w.closeCurrentFile()
						return 0, errors.Trace(err)
					}
				}
			} else {
				return 0, fmt.Errorf("Unable to seek to %s offset %v: file not long enough", filename, offset)
			}
		}

		o, err = w.currentFile.Seek(offset, os.SEEK_SET)
		if err != nil {
			w.closeCurrentFile()
			return 0, errors.Trace(err)
		}
	}

	return offset, nil
}

func (w *RawdataStorage) closeCurrentFile() (err error) {
	if w.currentFile != nil {
		err = w.currentFile.Close()
		w.currentFile = nil
		w.currentFilename = ""
	}
	return err
}
