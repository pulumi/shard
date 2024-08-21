package ignore

import "testing"

func TestMain(m *testing.M) {}

func Testament(t *testing.T) {}

func TeståMent(t *testing.T) {}

func testPrivate(t *testing.T) {}

func TestWithExtraArgs(t *testing.T, foo int) {}

func BenchmarkIgnored(b *testing.B) {}

func OtherTest(t *testing.T) {}

func ExampleWithNoOutput() {}

func TestWeirdArg(t *testing.Cover) {}
func TestWeirdArg2(t testing.T)     {}
