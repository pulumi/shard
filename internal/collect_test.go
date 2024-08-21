package internal

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCollect(t *testing.T) {

	t.Run("ignored", func(t *testing.T) {
		tests, err := Collect("testdata/ignored")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if len(tests) == 0 {
			return
		}
		t.Fatalf("unexepected tests collected: %+v", tests)
	})

	t.Run("collected", func(t *testing.T) {
		want := []testf{
			{Path: "testdata/collected", Name: "Test"},
			{Path: "testdata/collected", Name: "Test_It"},
			{Path: "testdata/collected", Name: "TestⱯUnicodeName"},
			{Path: "testdata/collected", Name: "TestAliasedT"},
			{Path: "testdata/collected", Name: "TestMain"},
			{Path: "testdata/collected", Name: "FuzzIt"},
			{Path: "testdata/collected", Name: "TestWithSubtests"},
			{Path: "testdata/collected", Name: "ExampleWithOutput"},
		}

		tests, err := Collect("testdata/collected")
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		if !reflect.DeepEqual(want, tests) {
			t.Fatalf("wanted\n\t%+v\nbut got\n\t%+v", want, tests)
		}
	})

}

func TestAssign(t *testing.T) {
	given := []testf{
		{Path: "package", Name: "Test"},
		{Path: "package", Name: "Test0"},
		{Path: "package", Name: "Test1"},
		{Path: "package", Name: "Test2"},
		{Path: "other/package", Name: "Test"},
		{Path: "other/package", Name: "Test3"},
		{Path: "other/package", Name: "Test4"},
		{Path: "other/package", Name: "Test5"},
		{Path: "other/package", Name: "Test6"},
	}

	tests := []struct {
		name         string
		index, total int
		seed         int64
		wantPaths    []string
		wantNames    []string
	}{
		{
			name:      "1/1",
			index:     0,
			total:     1,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test", "Test0", "Test1", "Test2", "Test3", "Test4", "Test5", "Test6"},
		},
		{
			index:     0,
			total:     100,
			wantPaths: []string{"./package"},
			wantNames: []string{"Test"},
		},
		{
			index:     99,
			total:     100,
			wantPaths: nil,
			wantNames: nil,
		},
		{
			index:     0,
			total:     4,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test", "Test6"},
		},
		{
			index:     1,
			total:     4,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test0", "Test3"},
		},
		{
			index:     2,
			total:     4,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test1", "Test4"},
		},
		{
			index:     3,
			total:     4,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test2", "Test5"},
		},
		{
			index:     0,
			total:     3,
			seed:      1,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test2", "Test3", "Test4"},
		},
		{
			index:     1,
			total:     3,
			seed:      1,
			wantPaths: []string{"./other/package"},
			wantNames: []string{"Test", "Test5", "Test6"},
		},
		{
			index:     2,
			total:     3,
			seed:      1,
			wantPaths: []string{"./other/package", "./package"},
			wantNames: []string{"Test", "Test1"},
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%d/%d", tt.index, tt.total), func(t *testing.T) {
			names, paths := Assign(given, tt.index, tt.total, tt.seed)
			if !reflect.DeepEqual(tt.wantPaths, paths) {
				t.Errorf("wanted\n\t%+v\nbut got\n\t%+v", tt.wantPaths, paths)
			}
			if !reflect.DeepEqual(tt.wantNames, names) {
				t.Errorf("wanted\n\t%+v\nbut got\n\t%+v", tt.wantNames, names)
			}
		})
	}
}
