package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	name            string
	compareFilePath string
	limit           int64
	offset          int64
	err             error
}

func TestCopy(t *testing.T) {
	fromFilePath := "testdata/input.txt"
	tests := []testCase{
		{
			name:            "offset 0 limit 0",
			compareFilePath: "testdata/out_offset0_limit0.txt",
			limit:           0,
			offset:          0,
		},
		{
			name:            "offset 0 limit 10",
			compareFilePath: "testdata/out_offset0_limit10.txt",
			offset:          0,
			limit:           10,
		},
		{
			name:            "offset 0 limit 1000",
			compareFilePath: "testdata/out_offset0_limit1000.txt",
			offset:          0,
			limit:           1000,
		},
		{
			name:            "offset 0 limit 10000",
			compareFilePath: "testdata/out_offset0_limit10000.txt",
			offset:          0,
			limit:           10000,
		},
		{
			name:            "offset 100 limit 1000",
			compareFilePath: "testdata/out_offset100_limit1000.txt",
			offset:          100,
			limit:           1000,
		},
		{
			name:            "offset 6000 limit 1000",
			compareFilePath: "testdata/out_offset6000_limit1000.txt",
			offset:          6000,
			limit:           1000,
		},
		{
			name:   "Error offset exceed file size",
			offset: 10000,
			limit:  0,
			err:    ErrOffsetExceedsFileSize,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			tempFile, err := os.CreateTemp("/tmp", "out")
			defer func(tempFile *os.File) {
				errClose := tempFile.Close()
				if errClose != nil {
					err = errClose
				}
			}(tempFile)

			if err != nil {
				fmt.Println("Error creating temp file:", err)
				return
			}

			err = Copy(fromFilePath, tempFile.Name(), test.offset, test.limit)
			if err != nil {
				assert.Equal(t, test.err, err)
			} else {
				outHash, _ := fileHash(t, tempFile.Name())
				compareHash, _ := fileHash(t, test.compareFilePath)
				assert.Equal(t, outHash, compareHash)
			}
		})
	}
}

func fileHash(t *testing.T, filename string) (string, error) {
	t.Helper()
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
