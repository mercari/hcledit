package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.mercari.io/hcledit"
)

type UpdateOptions struct {
	Type string
}

func NewCmdUpdate() *cobra.Command {
	opts := &UpdateOptions{}
	cmd := &cobra.Command{
		Use:   "update <query> <value> <file>",
		Short: "Update the given field with a value",
		Long:  `Runs an address query on a hcl file and update the given field with a value.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			return runUpdate(opts, args)
		},
	}

	cmd.Flags().StringVarP(&opts.Type, "type", "t", "string", "Type of the value")

	return cmd
}

func runUpdate(opts *UpdateOptions, args []string) error {
	query, valueStr, filePath := args[0], args[1], args[2]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	value, err := convert(valueStr, opts.Type)
	if err != nil {
		return fmt.Errorf("failed to convert input to specific type: %s", err)
	}

	if err := editor.Update(query, value); err != nil {
		return fmt.Errorf("failed to delete: %s", err)
	}

	return editor.OverWriteFile()
}
