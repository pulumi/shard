package main

import (
	"errors"
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
	output := flag.String("output", "", "output format (env)")

	flag.Parse()

	p := prog{index: *index, total: *total, seed: *seed, root: *root, output: *output}
	out, err := p.run()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Fprint(os.Stdout, out)
}

type prog struct {
	index  int
	total  int
	seed   int64
	root   string
	output string
}

func (p prog) run() (string, error) {
	if p.index < 0 {
		return "", errors.New("index is required")
	}
	if p.total < 0 {
		return "", errors.New("total is required")
	}
	if p.index >= p.total {
		return "", errors.New("index must be less than total")
	}

	tests, err := internal.Collect(p.root)
	if err != nil {
		log.Fatal(err)
	}

	names, paths := internal.Assign(tests, p.index, p.total, p.seed)

	// No-op if we didn't find any tests or get any assigned.
	if len(paths) == 0 {
		paths = []string{p.root}
		names = []string{"NoTestsFound"}
	}

	pattern := fmt.Sprintf(`^(?:%s)\$`, strings.Join(names, "|"))

	switch p.output {
	case "env":
		return fmt.Sprintf("SHARD_TESTS=%s\nSHARD_PATHS=%s", pattern, strings.Join(paths, " ")), nil
	default:
		return fmt.Sprintf("-run %s  %s", pattern, strings.Join(paths, " ")), nil
	}
}
