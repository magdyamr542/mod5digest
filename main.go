package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

func main() {
	// P()
	// C()
	CB()
}

// Parallel
func P() {
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

// Concurrent
func C() {
	if len(os.Args[1:]) > 1 {
		log.Fatal("only specify path to directory as argument.")
	}
	dir := os.Args[1]
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	resultc, err := ConcurrentMd5All(ctx, dir)
	if err != nil {
		log.Fatal(err)
	}

	pathMap := make(map[string][md5.Size]byte)
	for result := range resultc {
		if result.err != nil {
			log.Fatalf("calculate checksum concurrent: %v\n", result.err)
			os.Exit(1)
		}
		pathMap[result.file] = result.checksum
	}

	end := time.Since(start)

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

// Concurrent Bounded
func CB() {
	if len(os.Args[1:]) > 1 {
		log.Fatal("only specify path to directory as argument.")
	}
	dir := os.Args[1]
	start := time.Now()

	workers := 30
	done := make(chan struct{})
	resultc, errc := ConcurrentBoundedMd5All(done, dir, workers)

	go func() {
		<-time.After(time.Millisecond * 300)
		fmt.Println("call to done from main")
		close(done)
	}()

	pathMap := make(map[string][md5.Size]byte)
	for result := range resultc {
		if result.err != nil {
			log.Fatalf("calculate checksum concurrent: %v\n", result.err)
			os.Exit(1)
		}
		pathMap[result.file] = result.checksum
	}

	end := time.Since(start)
	fmt.Println("duration", end)

	if err := <-errc; err != nil {
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

}
