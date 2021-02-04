package command

import (
	"testing"
)

func TestRunRead(t *testing.T) {
	defaultOpts := &ReadOptions{
		OutputFormat: "go-template='{{.Value}}'",
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
			opts: &ReadOptions{
				OutputFormat: "go-template='prefix {{.Value}} suffix'",
			},
		},
		"formatted string with =": {
			query: "module.my-module.string_variable",
			want:  "module.my-module.string_variable=string",
			opts: &ReadOptions{
				OutputFormat: "go-template='{{.Key}}={{.Value}}'",
			},
		},
		"key and value string": {
			query: "module.my-module.string_variable",
			want:  "module.my-module.string_variable string",
			opts: &ReadOptions{
				OutputFormat: "go-template='{{.Key}} {{.Value}}'",
			},
		},
		"json string": {
			query: "module.my-module.string_variable",
			want:  `{"module.my-module.string_variable":"string"}`,
			opts: &ReadOptions{
				OutputFormat: "json",
			},
		},
		"yaml string": {
			query: "module.my-module.string_variable",
			want:  "module.my-module.string_variable: string\n",
			opts: &ReadOptions{
				OutputFormat: "yaml",
			},
		},
		"array": {
			query: "module.my-module.array_variable",
			want:  "[a b c]",
			opts:  defaultOpts,
		},
		"yaml array": {
			query: "module.my-module.array_variable",
			want: `module.my-module.array_variable:
- a
- b
- c
`,
			opts: &ReadOptions{
				OutputFormat: "yaml",
			},
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
		tc := tc
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
