package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// Md5All returns a map of a file to its md5 checksum and an error
func Md5All(directory string) (map[string][md5.Size]byte, error) {
	files, err := os.ReadDir(directory)
	checksumMap := make(map[string][md5.Size]byte)
	if err != nil {
		return map[string][16]byte{}, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		data, err := ioutil.ReadFile(filepath.Join(directory, file.Name()))
		if err != nil {
			return map[string][16]byte{}, err
		}

		checksum := md5.Sum(data)
		checksumMap[file.Name()] = checksum
	}
	return checksumMap, nil
}

func main() {
	if len(os.Args[1:]) > 1 {
		log.Fatal("only specify path to directory as argument.")
	}
	dir := os.Args[1]
	start := time.Now()
	pathMap, err := Md5All(dir)
	end := time.Since(start)

	if err != nil {
		log.Fatal(err)
	}
	var paths []string
	for path := range pathMap {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		fmt.Printf("%x %s\n", pathMap[path], path)
	}
	fmt.Println("duration", end)
}
