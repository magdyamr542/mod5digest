package main

import (
	"crypto/md5"
	"io/ioutil"
	"os"
	"path/filepath"
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
