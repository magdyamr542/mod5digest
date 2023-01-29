package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sync"
)

func emitPaths(done <-chan struct{}, directory string) (<-chan string, <-chan error) {
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
			case <-done:
				return errors.New("done")

			}
			return nil
		})
	}()
	return paths, errorc
}

func digester(done <-chan struct{}, id int, paths <-chan string, resultc chan<- result) {
	for path := range paths {
		fmt.Printf("[%d] process %s\n", id, path)
		data, err := ioutil.ReadFile(path)
		checksum := md5.Sum(data)
		select {
		case resultc <- result{err: err, file: path, checksum: checksum}:
		case <-done:
			fmt.Printf("digester %d done\n", id)
			return
		}
	}
}
func ConcurrentBoundedMd5All(done <-chan struct{}, directory string, workers int) (<-chan result, <-chan error) {
	paths, errorc := emitPaths(done, directory)
	resultc := make(chan result)

	var wg sync.WaitGroup

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			digester(done, id, paths, resultc)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		fmt.Printf("all digesters done\n")
		close(resultc)
	}()

	return resultc, errorc
}
