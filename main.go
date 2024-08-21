package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/blampe/shard/internal"
)

type testf struct {
	path string
	name string
}

func main() {
	log.SetFlags(0)
	log.SetPrefix(filepath.Base(os.Args[0]) + ": ")

	root := flag.String("root", ".", "directory to search for tests")
	index := flag.Int("index", -1, "shard index to collect tests for")
	total := flag.Int("total", -1, "total number of shards")
	seed := flag.Int64("seed", 0, "randomly shuffle tests using this seed")

	flag.Parse()
	if *index < 0 {
		log.Fatal("index is required")
	}
	if *total < 0 {
		log.Fatal("total is required")
	}
	if *index >= *total {
		log.Fatal("index must be less than total")
	}

	tests, err := internal.Collect(*root)
	if err != nil {
		log.Fatal(err)
	}

	names, paths := internal.Assign(tests, *index, *total, *seed)

	// No-op if we didn't find any tests or get any assigned.
	if len(paths) == 0 {
		paths = []string{*root}
		names = []string{"NoTestsFound"}
	}

	fmt.Printf("-run ^(%s)$ %s\n", strings.Join(names, "|"), strings.Join(paths, " "))
}
