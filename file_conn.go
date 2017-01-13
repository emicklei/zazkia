package main

import (
	"io"
	"log"
	"os"
)

type fileConnection struct {
	readFilename  string
	readable      io.ReadCloser
	writeable     io.WriteCloser
	writeFilename string
}

// Read implements io.Reader
func (f *fileConnection) Read(p []byte) (n int, err error) {
	if f.readable != nil {
		n, err = f.readable.Read(p)
		if io.EOF == err {
			return 0, nil
		}
	}
	return 0, nil
}

// Write implements io.Writer
func (f *fileConnection) Write(p []byte) (n int, err error) {
	return 0, nil
}

// Close implements io.Closer
func (f *fileConnection) Close() error {
	if f.readable != nil {
		f.readable.Close()
	}
	return nil
}

func (f *fileConnection) String() string {
	return "file"
}

func newFileConnection(src string) *fileConnection {
	r, err := os.Open(src)
	if err != nil {
		log.Printf("unable to open for read:%s error:%v", src, err)
	}
	return &fileConnection{
		readFilename: src,
		readable:     r,
	}
}
