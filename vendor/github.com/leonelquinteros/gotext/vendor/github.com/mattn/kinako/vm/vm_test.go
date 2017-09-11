package vm

import (
	"reflect"
	"testing"
)

func TestExecute(t *testing.T) {
	e := NewEnv()
	e.Define("foo", int64(1))
	e.Define("bar", int64(2))
	e.Define("baz", int64(3))

	tests := []struct {
		input string
		want  interface{}
	}{
		{
			input: "foo+bar",
			want:  int64(3),
		},
		{
			input: "foo-bar",
			want:  int64(-1),
		},
		{
			input: "foo*bar",
			want:  int64(2),
		},
		{
			input: "foo/bar",
			want:  float64(0.5),
		},
		{
			input: "baz*(foo+bar)",
			want:  int64(9),
		},
		{
			input: "baz > 2 ? foo : bar",
			want:  int64(1),
		},
	}

	for _, tt := range tests {
		r, err := e.Execute(tt.input)
		if err != nil {
			t.Fatal(err)
		}
		got := r.Interface()
		if !reflect.DeepEqual(got, tt.want) {
			t.Fatalf("want %v, but %v:", tt.want, got)
		}
	}
}
