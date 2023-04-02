package httpwriter

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

const (
	fileNameMaxLength = 200
)

func NewDirectoryWriter(directory *Directory) (func(*http.Request) (io.WriteCloser, error), error) {
	if directory.Path == "" {
		return nil, errors.New("directory.Path must not be empty")
	}

	if err := os.MkdirAll(directory.Path, 0744); err != nil {
		return nil, err
	}

	if directory.CreateFileName == nil {
		directory.CreateFileName = func(req *http.Request, counter uint64) (string, error) {
			fileName := fmt.Sprintf(
				"%04d_%s@%s",
				counter,
				req.Method,
				strings.ReplaceAll(req.URL.String(), "/", "@"),
			)
			if len(fileName) > fileNameMaxLength {
				return fileName[0:fileNameMaxLength], nil
			}
			return fileName, nil
		}
	}

	return func(req *http.Request) (io.WriteCloser, error) {
		filePath, err := directory.CreateFilePath(req)
		if err != nil {
			return nil, err
		}
		return os.Create(filePath)
	}, nil
}

func MustDirectoryWriter(directory *Directory) func(*http.Request) (io.WriteCloser, error) {
	getWriter, err := NewDirectoryWriter(directory)
	if err != nil {
		panic(err)
	}
	return getWriter
}

type Directory struct {
	Path           string
	CreateFileName func(req *http.Request, counter uint64) (string, error)

	counter uint64
}

func (d *Directory) CreateFilePath(req *http.Request) (string, error) {
	counter := atomic.AddUint64(&d.counter, 1)
	fileName, err := d.CreateFileName(req, counter)
	if err != nil {
		return "", err
	}
	return filepath.Join(d.Path, fileName), nil
}
