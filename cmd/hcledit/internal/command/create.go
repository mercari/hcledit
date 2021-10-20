package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.mercari.io/hcledit"
)

type CreateOptions struct {
	Type    string
	After   string
	Comment string
}

func NewCmdCreate() *cobra.Command {
	opts := &CreateOptions{}
	cmd := &cobra.Command{
		Use:   "create <query> <value> <file>",
		Short: "Create a new field",
		Long:  `Runs an address query on a hcl file and create new field with given value.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			return runCreate(opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "string", "Type of the value")
	cmd.Flags().StringVarP(&opts.After, "after", "a", "", "Field key which before the value will be created")
	cmd.Flags().StringVarP(&opts.Comment, "comment", "c", "", "Comment to be inserted before the field added. Comment symbols like // are required")

	return cmd
}

func runCreate(opts *CreateOptions, args []string) error {
	query, valueStr, filePath := args[0], args[1], args[2]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	value, err := convert(valueStr, opts.Type)
	if err != nil {
		return fmt.Errorf("failed to convert input to specific type: %s", err)
	}

	if err := editor.Create(query, value, hcledit.WithAfter(opts.After), hcledit.WithComment(opts.Comment)); err != nil {
		return fmt.Errorf("failed to create: %s", err)
	}

	return editor.OverWriteFile()
}
