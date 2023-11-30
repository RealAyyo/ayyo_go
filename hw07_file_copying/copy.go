package main

import (
	"errors"
	"io"
	"math"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromFilePath, toFilePath string, offset, limit int64) error {
	fileInfo, err := os.Stat(fromFilePath)
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 {
		limit = fileInfo.Size()
	}

	bufSize := 1024

	if int(limit) < bufSize {
		bufSize = int(limit)
	}

	buf := make([]byte, bufSize)

	fileFrom, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer func() {
		errClose := fileFrom.Close()
		if errClose != nil {
			err = errClose
		}
	}()

	fileTo, err := os.Create(toFilePath)
	if err != nil {
		return err
	}
	defer func() {
		errClose := fileTo.Close()
		if errClose != nil {
			err = errClose
		}
	}()

	if _, err = fileFrom.Seek(offset, io.SeekStart); err != nil {
		return err
	}

	copyProcess(fileFrom, fileTo, buf, limit, bufSize)

	return err
}

func copyProcess(in io.Reader, out io.Writer, buf []byte, limit int64, bufSize int) {
	var written int

	iterationCount := math.Ceil(float64(limit) / float64(bufSize))
	bar := progressbar.Default(int64(iterationCount))

	for {
		bar.Add(1)
		if written >= int(limit) {
			break
		}

		n, err := in.Read(buf)
		if err != nil && err != io.EOF {
			break
		}

		if n == 0 {
			break
		}

		bytesLeft := int(limit) - written
		if n > bytesLeft {
			n = bytesLeft
		}

		if _, err := out.Write(buf[:n]); err != nil {
			break
		}

		written += n
	}
}
