package main

import (
	"errors"
	"testing"
)

func TestProg(t *testing.T) {
	tests := []struct {
		p       prog
		name    string
		want    string
		wantErr error
	}{
		{
			name: "default output",
			p:    prog{total: 1, root: "."},
			want: `-run "^(?:TestAssign|TestCollect|TestProg)\$"  ./. ./internal`,
		},
		{
			name: "env output",
			p:    prog{output: "env", total: 1, root: "."},
			want: `SHARD_TESTS="^(?:TestAssign|TestCollect|TestProg)\$"
SHARD_PATHS="./. ./internal"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := tt.p.run()
			if out != tt.want {
				t.Errorf("wanted %q but got %q", tt.want, out)

			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted %q but got %q", tt.wantErr, err)
			}
		})
	}

}
