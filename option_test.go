package hcledit_test

import (
	"strings"
	"testing"

	"go.mercari.io/hcledit"
)

func TestWithComment(t *testing.T) {
	cases := map[string]struct {
		input string
		exec  func(editor *hcledit.HCLEditor) error
		want  string
	}{
		"CreateAttribute": {
			input: `
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("attribute", "str", hcledit.WithComment("// Comment"))
			},
			want: `
// Comment
attribute = "str"
`,
		},

		"CreateAttributeInBlock": {
			input: `
block "label1" {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("block.label1.attribute", "str", hcledit.WithComment("// Comment"))
			},
			want: `
block "label1" {
  // Comment
  attribute = "str"
}
`,
		},

		"CreateAttributeInObject": {
			input: `
object = {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("object.attribute", "str", hcledit.WithComment("// Comment"))
			},
			want: `
object = {
  // Comment
  attribute = "str"
}
`,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			editor, err := hcledit.Read(strings.NewReader(tc.input), "")
			if err != nil {
				t.Fatal(err)
			}

			if err := tc.exec(editor); err != nil {
				t.Fatal(err)
			}

			got := string(editor.Bytes())
			if got != tc.want {
				t.Errorf("Update() mismatch:\ngot:%s\nwant:%s\n", got, tc.want)
			}
		})
	}
}

func TestWithAfter(t *testing.T) {
	cases := map[string]struct {
		input string
		query string
		after string
		want  string
	}{
		"CreateAttribute": {
			input: `
attribute1 = "str1"
attribute2 = "str2"
`,
			query: "attribute3",
			after: "attribute1",
			want: `
attribute1 = "str1"
attribute3 = "C"
attribute2 = "str2"
`,
		},

		"CreateAttributeInBlock": {
			input: `
block {
  attribute1 = "str1"
  attribute2 = "str2"
}
`,
			query: "block.attribute3",
			after: "attribute1",
			want: `
block {
  attribute1 = "str1"
  attribute3 = "C"
  attribute2 = "str2"
}
`,
		},

		"CreateAttributeInObject": {
			input: `
object = {
  attribute1 = "str1"
  attribute2 = "str2"
}
`,
			query: "object.attribute3",
			after: "attribute1",
			want: `
object = {
  attribute1 = "str1"
  attribute3 = "C"
  attribute2 = "str2"
}
`,
		},

		"CreateAttributeWithBlock": {
			input: `
attribute1 = "str1"
block {
  attribute2 = "str2"
}
`,
			query: "attribute3",
			after: "attribute1",
			want: `
attribute1 = "str1"
attribute3 = "C"
block {
  attribute2 = "str2"
}
`,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			editor, err := hcledit.Read(strings.NewReader(tc.input), "")
			if err != nil {
				t.Fatal(err)
			}

			if err := editor.Create(tc.query, "C", hcledit.WithAfter(tc.after)); err != nil {
				t.Fatal(err)
			}

			got := string(editor.Bytes())
			if got != tc.want {
				t.Errorf("Update() mismatch:\ngot:%s\nwant:%s\n", got, tc.want)
			}
		})
	}
}

func TestWithNewLine(t *testing.T) {
	cases := map[string]struct {
		input string
		exec  func(editor *hcledit.HCLEditor) error
		want  string
	}{
		"CreateAttribute": {
			input: `
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("attribute", "str", hcledit.WithNewLine())
			},
			want: `

attribute = "str"
`,
		},

		"CreateAttributeInBlock": {
			input: `
block "label1" {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("block.label1.attribute", "str", hcledit.WithNewLine())
			},
			want: `
block "label1" {

  attribute = "str"
}
`,
		},

		"CreateAttributeInObject": {
			input: `
object = {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("object.attribute", "str", hcledit.WithNewLine())
			},
			want: `
object = {

  attribute = "str"
}
`,
		},
		"CreateWithComment": {
			input: `
object = {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("object.attribute", "str", hcledit.WithComment("// Comment"), hcledit.WithNewLine())
			},
			want: `
object = {

  // Comment
  attribute = "str"
}
`,
		},

		"CreateWithQuerySeparator": {
			input: `
object = {
}
`,
			exec: func(editor *hcledit.HCLEditor) error {
				return editor.Create("object/attribute", "str", hcledit.WithQuerySeparator('/'))
			},
			want: `
object = {
  attribute = "str"
}
`,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			editor, err := hcledit.Read(strings.NewReader(tc.input), "")
			if err != nil {
				t.Fatal(err)
			}

			if err := tc.exec(editor); err != nil {
				t.Fatal(err)
			}

			got := string(editor.Bytes())
			if got != tc.want {
				t.Errorf("Update() mismatch:\ngot:%s\nwant:%s\n", got, tc.want)
			}
		})
	}
}
