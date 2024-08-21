package collected

import (
	"fmt"
	"testing"
)

type T = testing.T

func Test(t *testing.T) {}

func Test_It(t *testing.T) {}

func TestⱯUnicodeName(t *testing.T) {}

func TestAliasedT(t *T) {}

func TestMain(t *testing.T) {} // Only testing.M is ignored.

func FuzzIt(f *testing.F) {}

func ExampleWithOutput() {
	fmt.Println("foo")
	// Output: foo
}

func TestWithSubtests(t *testing.T) {
	t.Run("subtest", func(t *testing.T) {})
}
