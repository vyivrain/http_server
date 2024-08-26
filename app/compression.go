package main

type compressionMethod interface {
	dataContainsCompressionMethod(data []byte) bool
	compress(data []byte) ([]byte, error)
	decompress(data []byte) ([]byte, error)
	name() string
}
