package main

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func emitPaths(ctx context.Context, directory string) (<-chan string, <-chan error) {
	paths := make(chan string)
	errorc := make(chan error, 1)
	go func() {
		defer close(paths)
		errorc <- filepath.Walk(directory, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			select {
			case paths <- path:
			case <-ctx.Done():
				return errors.New("context canceled")

			}
			return nil
		})
	}()
	return paths, errorc
}

func digester(ctx context.Context, id int, paths <-chan string, resultc chan<- result) {
	for path := range paths {
		fmt.Printf("[%d] process %s\n", id, path)
		data, err := ioutil.ReadFile(path)
		checksum := md5.Sum(data)
		select {
		case resultc <- result{err: err, file: path, checksum: checksum}:
		case <-ctx.Done():
			return
		}
	}
}
func ConcurrentBoundedMd5All(ctx context.Context, directory string, workers int) (<-chan result, <-chan error) {
	paths, errorc := emitPaths(ctx, directory)
	resultc := make(chan result)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			digester(ctx, id, paths, resultc)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(resultc)
	}()

	return resultc, errorc
}
