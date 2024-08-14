package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var re = regexp.MustCompile(`^Test[A-Z_]`)

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

	tests := []testf{}

	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file to find test functions
		fileSet := token.NewFileSet()
		node, err := parser.ParseFile(fileSet, path, nil, 0)
		if err != nil {
			return err
		}
		for _, decl := range node.Decls {
			f, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			name := f.Name.Name
			if !re.MatchString(name) {
				continue
			}
			tests = append(tests, testf{path: filepath.Dir(path), name: name})
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Shuffle the tests.
	if *seed != 0 {
		random := rand.New(rand.NewSource(*seed))
		for i := range tests {
			j := random.Intn(i + 1) //nolint:gosec // Not cryptographic.
			tests[i], tests[j] = tests[j], tests[i]
		}
	}

	// Assign tests to our shard.
	paths := []string{}
	names := []string{}
	for idx, test := range tests {
		if idx%*total != *index {
			continue
		}
		paths = append(paths, "./"+test.path)
		names = append(names, test.name)
	}

	// De-dupe.
	slices.Sort(paths)
	slices.Sort(names)
	paths = slices.Compact(paths)
	names = slices.Compact(names)

	// No-op if we didn't find any tests or get any assigned.
	if len(paths) == 0 {
		paths = []string{*root}
		names = []string{"NoTestsFound"}
	}

	fmt.Printf("-run ^(%s)$ %s\n", strings.Join(names, "|"), strings.Join(paths, " "))
}
