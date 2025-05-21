package internal

import (
	"fmt"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"math/rand"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

type testf struct {
	Path string
	Name string
}

func isGoModule(path string) bool {
	_, err := os.Lstat(filepath.Join(path, "go.mod"))
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		panic(err)
	}
	return true
}

func Collect(root string) ([]testf, error) {
	tests := []testf{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't collect testdata directories.
		if info.IsDir() &&
			(filepath.Base(info.Name()) == "testdata" || filepath.Base(info.Name()) == "vendor" ||
				(info.Name() != root && isGoModule(info.Name()))) {
			return filepath.SkipDir
		}

		if info.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file to find test functions
		file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}
			if !isTestFunc(fn) {
				continue
			}
			tests = append(tests, testf{Path: filepath.Dir(path), Name: fn.Name.Name})
		}

		examples := doc.Examples(file)
		for _, ex := range examples {
			if ex.Output == "" && !ex.EmptyOutput {
				// Don't run tests with empty output.
				continue
			}
			tests = append(tests, testf{Path: filepath.Dir(path), Name: "Example" + ex.Name})

		}

		return nil
	})
	return tests, err
}

func CheckForClashingTestNames(tests []testf) []string {
	slices.SortStableFunc(tests, func(a, b testf) int {
		return strings.Compare(a.Name, b.Name)
	})

	repeats := map[string]struct{}{}
	prefixed := map[string]int{}

	for i := 0; i < len(tests)-1; i++ {
		for j := i + 1; j < len(tests); j++ {
			if tests[i].Name == tests[j].Name {
				repeats[tests[i].Name] = struct{}{}
				continue
			}
			if strings.HasPrefix(tests[i].Name, tests[j].Name) {
				prefixed[tests[j].Name]++
			}
			if strings.HasPrefix(tests[j].Name, tests[i].Name) {
				prefixed[tests[i].Name]++
			}
		}
	}

	warnings := []string{}
	for name := range repeats {
		warnings = append(warnings, fmt.Sprintf("test name %q exists in multiple packages, consider renaming it", name))
	}
	for name := range prefixed {
		warnings = append(warnings, fmt.Sprintf("test name %q is prefixed by %d other tests, consider renaming it", name, prefixed[name]))
	}
	return warnings
}

func Assign(tests []testf, index int, total int, seed int64) (names, paths []string) {
	// Shuffle the tests.
	if seed != 0 {
		random := rand.New(rand.NewSource(seed))
		for i := range tests {
			j := random.Intn(i + 1) //nolint:gosec // Not cryptographic.
			tests[i], tests[j] = tests[j], tests[i]
		}
	}

	// Assign tests to our shard.
	for idx, test := range tests {
		if idx%total != index {
			continue
		}
		paths = append(paths, "./"+test.Path)
		names = append(names, test.Name)
	}

	// De-dupe.
	slices.Sort(names)
	slices.Sort(paths)
	paths = slices.Compact(paths)

	return names, paths
}

// isTestFunc tells whether fn has the type of a testing function. arg
// specifies the parameter type we look for: B, F, M or T.
func isTestFunc(fn *ast.FuncDecl) bool {
	if fn.Type.Results != nil && len(fn.Type.Results.List) > 0 ||
		fn.Type.Params.List == nil ||
		len(fn.Type.Params.List) != 1 ||
		len(fn.Type.Params.List[0].Names) > 1 {
		return false
	}
	ptr, ok := fn.Type.Params.List[0].Type.(*ast.StarExpr)
	if !ok {
		return false
	}

	fnName := fn.Name.Name

	if !(strings.HasPrefix(fnName, "Test") || strings.HasPrefix(fnName, "Fuzz")) {
		return false
	}

	if len(fnName) > 4 {
		rune, _ := utf8.DecodeRuneInString(fnName[4:])
		if unicode.IsLower(rune) {
			return false
		}
	}

	for _, want := range []string{"T", "F"} {
		// We can't easily check that the type is *testing.M
		// because we don't know how testing has been imported,
		// but at least check that it's *M or *something.M.
		// Same applies for B, F and T.
		if name, ok := ptr.X.(*ast.Ident); ok && name.Name == want {
			return true
		}
		if sel, ok := ptr.X.(*ast.SelectorExpr); ok && sel.Sel.Name == want {
			return true
		}
	}
	return false
}
