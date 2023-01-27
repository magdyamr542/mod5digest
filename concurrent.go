package main

import (
	"context"
	"crypto/md5"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

type result struct {
	err      error
	file     string
	checksum [md5.Size]byte
}

func ConcurrentMd5All(ctx context.Context, directory string) (<-chan result, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	resultc := make(chan result)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		select {
		case <-ctx.Done():
			return resultc, errors.New("context canceled")
		default:
		}

		wg.Add(1)
		go func(fileName string) {
			defer wg.Done()
			data, err := ioutil.ReadFile(filepath.Join(directory, fileName))
			checksum := md5.Sum(data)
			select {
			case <-ctx.Done():
				return
			case resultc <- result{err: err, file: fileName, checksum: checksum}:
			}
		}(file.Name())
	}

	// close the channel
	go func() {
		wg.Wait()
		close(resultc)
	}()

	return resultc, nil
}
