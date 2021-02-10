package command

import (
	"fmt"

	"github.com/spf13/cobra"

	"go.mercari.io/hcledit"
)

func NewCmdDelete() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <query> <file>",
		Short: "Delete the given field and its value",
		Long:  `Runs an address query on a hcl file and delete the given field and its value.`,
		Args:  cobra.ExactArgs(2),
		RunE: func(_ *cobra.Command, args []string) error {
			return runDelete(args)
		},
	}
}

func runDelete(args []string) error {
	query, filePath := args[0], args[1]

	editor, err := hcledit.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	if err := editor.Delete(query); err != nil {
		return fmt.Errorf("failed to delete: %s", err)
	}

	return editor.OverWriteFile()
}
