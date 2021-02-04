package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.mercari.io/hcledit"
)

func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update <query> <value> <file>",
		Short: "Update the given field with a value",
		Long:  `Runs an address query on a hcl file and update the given field with a value.`,
		Args:  cobra.ExactArgs(3),
		RunE: func(_ *cobra.Command, args []string) error {
			if err := runUpdate(args); err != nil {
				return err
			}
			return nil
		},
	}
	return cmd
}

func runUpdate(args []string) error {
	query := args[0]
	value := args[1]
	filePath := args[2]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if err := editor.Update(query, value); err != nil {
		return fmt.Errorf("failed to delete: %s", err)
	}

	return editor.OverWriteFile()
}
