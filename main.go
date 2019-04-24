package main

import (
	"github.com/oct2pus/emoji-dl/emoji"
	"flag"
	"fmt"
	"os"
	"strings"
)

var (
	flagVerbose      bool
	flagSuperVerbose bool
	flagBatch        bool
	flagSize         int
)

func init() {
	flag.BoolVar(&flagVerbose,
		"v",
		false,
		"print successful downloads")
	flag.BoolVar(&flagSuperVerbose,
		"verbose",
		false,
		"print successful downloads and errors")
	flag.BoolVar(&flagBatch, "batch", true, "perform downloads in batches")
	flag.IntVar(&flagSize, "size", 75, "how many files to download at once, "+
		"ignored if batch=false")
}

func main() {
	flag.Parse()
	url := flag.Arg(0)

	if url == "" {
		fmt.Printf("Please enter a URL.\n")
		return
	}

	url, dir := processNames(url)

	collection, err := emoji.NewCollection(url)
	if has(err) { // fail if can't create Collection
		return
	}

	wd, err := os.Getwd()
	if has(err) { // fail if can't find working directory
		return
	}

	dir = wd + "/" + dir
	err = os.MkdirAll(dir, os.ModePerm)
	if has(err) { // fail if unable to create directory
		return
	}

	emoji.Verbose.Downloads = flagVerbose
	if flagSuperVerbose {
		emoji.Verbose.Downloads = flagSuperVerbose
		emoji.Verbose.Errors = flagSuperVerbose
	}

	emoji.Batch.Downloads = flagBatch
	emoji.Batch.Size = flagSize

	if flagSize <= 0 {
		emoji.Batch.Downloads = false
		emoji.Batch.Size = 1
	}

	go collection.DownloadAll(dir)

	for !collection.DownloadFinished() {
		if !emoji.Verbose.Downloads && !emoji.Verbose.Errors {
			fmt.Printf("\rdownloading...")
		}
	}

	fmt.Println("Finished!")
}

// processNames converts the following into acceptable outputs for paths
// and file names.
func processNames(arg string) (string, string) {
	// add https:// if arg doesn't have it
	arg2 := arg
	if !strings.HasPrefix(arg, "https://") &&
		!strings.HasPrefix(arg, "http://") {
		arg = "https://" + arg
	} else if strings.HasPrefix(arg, "http://") {
		// change http:// to https:// in arg
		arg2 = strings.TrimPrefix(arg2, "http://")
		arg = strings.Replace(arg, "http://", "https://", 1)
	} else {
		// trim prefix in arg2 to avoid illegal file names
		arg2 = strings.TrimPrefix(arg2, "https://")
	}

	return arg, arg2
}

// hasError is a helper function to print errors.
func has(err error) bool {
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return true
	}
	return false
}
