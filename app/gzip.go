package main

import (
	"bytes"
	"compress/gzip"
	"io"
)

type GzipCompression struct{}

func (gc *GzipCompression) dataContainsCompressionMethod(data []byte) bool {
	// Check if the data is long enough to contain the gzip file format
	if len(data) < 2 {
		return false
	}

	// Check for the gzip file format
	if data[0] != 0x1f || data[1] != 0x8b {
		return false
	}

	return true
}

func (gc *GzipCompression) compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)
	defer gzipWriter.Close()

	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (gc *GzipCompression) decompress(data []byte) ([]byte, error) {
	reader := bytes.NewReader(data)
	gzReader, err := gzip.NewReader(reader)
	if err != nil {
		return nil, err
	}
	output, err2 := io.ReadAll(gzReader)
	if err2 != nil {
		return nil, err2
	}

	return output, nil
}
