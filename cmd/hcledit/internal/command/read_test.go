package command

import (
	"testing"
)

func TestRunRead(t *testing.T) {
	defaultOpts := &ReadOptions{
		ValueFormat: "%v",
		ValueOnly:   true,
	}
	fixture := "fixture/file.tf"

	cases := map[string]struct {
		query string
		want  string
		opts  *ReadOptions
	}{
		"bool": {
			query: "module.my-module.bool_variable",
			want:  "true",
			opts:  defaultOpts,
		},
		"int": {
			query: "module.my-module.int_variable",
			want:  "1",
			opts:  defaultOpts,
		},
		"string": {
			query: "module.my-module.string_variable",
			want:  "string",
			opts:  defaultOpts,
		},
		"formatted string": {
			query: "module.my-module.string_variable",
			want:  "prefix string suffix",
			opts:  &ReadOptions{
				ValueFormat: "prefix %v suffix",
				ValueOnly: true,
			},
		},
		"key and value string": {
			query: "module.my-module.string_variable",
			want:  "module.my-module.string_variable string\n",
			opts:  &ReadOptions{
				ValueFormat: "%v",
				ValueOnly: false,
			},
		},
		"array": {
			query: "module.my-module.array_variable",
			want:  "[a b c]",
			opts:  defaultOpts,
		},
		"map": {
			query: "module.my-module.map_variable.string_variable",
			want:  "string",
			opts:  defaultOpts,
		},
		"does not exist": {
			query: "module.my-module.does.not.exist",
			want:  "",
			opts:  defaultOpts,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			args := []string{tc.query, fixture}
			got, err := runRead(tc.opts, args)

			if err != nil {
				t.Fatalf("unexpected err %s", err)
			}
			if got != tc.want {
				t.Errorf("got: %s, want %s", got, tc.want)
			}
		})
	}
}
